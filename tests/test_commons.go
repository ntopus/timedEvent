package tests

import (
	"bytes"
	"context"
	"devgit.kf.com.br/core/lib-queue/queue"
	"devgit.kf.com.br/core/lib-queue/queue_repository"
	"encoding/json"
	"fmt"
	"github.com/ivanmeca/timedEvent/application"
	"github.com/ivanmeca/timedEvent/application/modules/config"
	"github.com/ivanmeca/timedEvent/application/modules/database/collection_managment"
	"github.com/ivanmeca/timedEvent/application/modules/database/data_types"
	"github.com/ivanmeca/timedEvent/application/modules/routes"
	"github.com/onsi/gomega"
	"github.com/streadway/amqp"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"time"
)

const (
	APP_NAME           = "timed-event"
	DATE_FORMAT        = "2006-01-02 15:04:05Z"
	TEST_ENDPOINT      = "http://localhost:9010/v1/event"
	CONTENT_TYPE       = "Content-Type"
	CONTENT_TYPE_CE    = "application/cloudevents"
	CONTENT_TYPE_JSON  = "application/json"
	PUBLISH_DATE       = "publishDate"
	PUBLISH_QUEUE      = "publishQueue"
	PUBLISH_TYPE       = "publishtype"
	TEST_PUBLISH_QUEUE = "throwAt"
	TEST_PUBLISH_TYPE  = "dataOnly"
)

type MockEvent struct {
	SpecVersion  string      `json:"specversion"`
	Type         string      `json:"type"`
	Source       string      `json:"source"`
	ID           string      `json:"id"`
	PublishDate  string      `json:"publishdate"`
	PublishQueue string      `json:"publishqueue"`
	PublishType  string      `json:"publishtype"`
	Data         interface{} `json:"data"`
}

type MockData struct {
	Ref         string
	PublishDate string
}

type fnConsume func(queueName string, msg []byte, counter int) bool

func BuildApplication() {
	cwd, err := os.Getwd()
	gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
	fmt.Println(cwd)
	os.Chdir(cwd)
	command := exec.Command("make", "build-native-production")
	err = command.Run()
	gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
}

func GetBinPath() string {
	cwd, err := os.Getwd()
	gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
	return filepath.Join(cwd, "bin")
}

func GetConfigPath() string {
	return filepath.Join(GetBinPath(), "config.json")
}

func GetQueue(queueName string, threadLimit int) *queue.Queue {
	params := queue_repository.NewQueueRepositoryParams("randomUser", "randomPass", "srvqueue.module.ntopus.com.br", 5672)
	params.SetVHost("/timed")
	qr, err := queue_repository.NewQueueRepository(params)
	gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
	qParam := queue.NewQueueParams(queueName)
	qParam.SetThreadLimit(threadLimit)
	q, err := qr.QueueDeclare(qParam, false)
	gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
	return q
}

func InitQueue(queueName string, counter *int, consume fnConsume) *queue.Queue {
	mu := sync.Mutex{}
	mu.Lock()
	*counter = 0
	mu.Unlock()
	q := GetQueue(queueName, 5000)
	err := q.StartConsume(func(queueName string, msg []byte) bool {
		mu.Lock()
		defer mu.Unlock()
		*counter++
		return consume(queueName, msg, *counter)
	})
	gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
	return q
}

func ParseResp(resp *http.Response, dataContainer interface{}) *routes.JsendMessage {
	buf, err := ioutil.ReadAll(resp.Body)
	gomega.Expect(err).To(gomega.BeNil())
	respBody := routes.JsendMessage{}
	respBody.SetData(dataContainer)
	err = json.Unmarshal(buf, &respBody)
	gomega.Expect(err).To(gomega.BeNil())
	return &respBody
}

func PurgeQueue(queue string) {
	conn, err := amqp.Dial("amqp://randomUser:randomPass@srvqueue.module.ntopus.com.br:5672/timed")
	gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
	ch, err := conn.Channel()
	gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
	_, err = ch.QueueInspect(queue)
	gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
	_, err = ch.QueuePurge(queue, false)
	gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
}

func GetMockReader(mockData interface{}) (io.Reader, error) {
	data, err := json.Marshal(mockData)
	if err != nil {
		return nil, err
	}
	buff := bytes.NewReader(data)
	return buff, nil
}

/*
func RunApp() *gexec.Session {
	appPath := filepath.Join(GetBinPath(), APP_NAME)
	command := exec.Command(appPath, "-c="+GetConfigPath())
	session, err := gexec.Start(command, ginkgo.GinkgoWriter, ginkgo.GinkgoWriter)
	gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
	time.Sleep(400 * time.Millisecond)
	fmt.Println("Application is running")
	return session
}
/*/
func RunApp() context.Context {
	ctx := context.Background()
	appMan := application.NewApplicationManager(GetConfigPath())
	err := appMan.RunApplication(ctx)
	gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
	defer func() {
		fmt.Println("Killing application")
		ctx.Done()
	}()
	time.Sleep(500 * time.Millisecond)
	fmt.Println("Application is running")
	return ctx
} //*/

func SaveConfigFile() {
	err := config.ConfigSample(GetConfigPath())
	gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
}

func SendGetRequest(url string) (resp *http.Response, err error) {
	Headers := make(map[string]string)
	return SendGetRequestWithHeaders(url, Headers)
}

func SendGetRequestWithHeaders(url string, headers map[string]string) (resp *http.Response, err error) {
	client := http.Client{}
	req, err := http.NewRequest(http.MethodGet, url, nil)
	gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
	for key, value := range headers {
		req.Header.Set(key, value)
	}
	return client.Do(req)
}

func SendPostRequest(url string, body io.Reader) (resp *http.Response, err error) {
	Headers := make(map[string]string)
	Headers["Content-Type"] = "application/json"
	return SendPostRequestWithHeaders(url, body, Headers)
}

func SendPostRequestWithHeaders(url string, body io.Reader, headers map[string]string) (resp *http.Response, err error) {
	client := http.Client{}
	req, err := http.NewRequest(http.MethodPost, url, body)
	gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
	for key, value := range headers {
		req.Header.Set(key, value)
	}
	return client.Do(req)
}

func GetMockEvent(publihsDate time.Time, publishType string, ref string) MockEvent {
	return MockEvent{
		SpecVersion:  "0.2",
		Type:         "TestEvent",
		Source:       "sourceEvent",
		ID:           fmt.Sprintf("mockEvent%s", ref),
		PublishDate:  publihsDate.Format(DATE_FORMAT),
		PublishQueue: TEST_PUBLISH_QUEUE,
		PublishType:  publishType,
		Data:         MockData{Ref: ref, PublishDate: publihsDate.Format(DATE_FORMAT)},
	}
}

func ClearDB() {
	data, err := collection_managment.NewEventCollection().Read(nil)
	gomega.Expect(err).ShouldNot(gomega.HaveOccurred())

	for _, item := range data {
		collection_managment.NewEventCollection().DeleteItem([]string{item.ArangoId})
	}
}

func ReadDocument(id string) *data_types.ArangoCloudEvent {
	data, err := collection_managment.NewEventCollection().ReadItem(id)
	gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
	return data
}
