package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/ivanmeca/timedEvent/application/modules/config"
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

const APP_NAME = "timed-event"

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

func GetMockReader(mockData interface{}) (io.Reader, error) {
	data, err := json.Marshal(mockData)
	if err != nil {
		return nil, err
	}
	buff := bytes.NewReader(data)
	return buff, nil
}

func RunApp() *gexec.Session {
	appPath := filepath.Join(GetBinPath(), APP_NAME)
	command := exec.Command(appPath, "-c="+GetConfigPath())
	session, err := gexec.Start(command, ginkgo.GinkgoWriter, ginkgo.GinkgoWriter)
	gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
	time.Sleep(400 * time.Millisecond)
	fmt.Println("Application is running")
	return session
}

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

//
//
//type singleton struct {
//	Date time.Time
//}
//
//var instance *singleton
//var once sync.Once
//
//func getDate() time.Time {
//	once.Do(func() {
//		instance = &singleton{Date: time.Now()}
//	})
//	return instance.Date
//}
//
//func getMockMsg() Message {
//	msg := Message{}
//	msg.MessageID = "1234234"
//	msg.ParentMessageID = "4565"
//	msg.Source = "Teste source"
//	msg.Severity = 10
//	msg.Err = nil
//	msg.DateTime = getDate()
//	msg.MsgType = 100
//	msg.Version = "1"
//	msg.Payload = "\"data\":{ \"eventDate\": { \"date\": \"2018-07-09 19:17:36.000000\",\"timezone_type\": 3,\"timezone\": \"UTC\"}"
//	return msg
//}
//
//func getQueue(threadLimit int) *queue.Queue {
//	qr, err := queue_repository.NewQueueRepository(queue_repository.NewQueueRepositoryParams("randomUser", "randomPass", "srvqueue.module.ntopus.com.br", 5672))
//	params := queue.NewQueueParams("newQueue")
//	params.SetThreadLimit(threadLimit)
//	q, err := qr.QueueDeclare(params, true)
//	gomega.Expect(err).To(gomega.BeNil())
//	return q
//}
//
//func publishMockMessageOnQueue(q *queue.Queue) {
//	msg := getMockMsg()
//	err := q.Publish(msg)
//	gomega.Expect(err).To(gomega.BeNil())
//}
//
//func getTcpConnection(address string) net.Conn {
//	var err error
//	connection, err = net.Dial("tcp", address)
//	gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
//	return connection
//}
//
//func tcpDataSender(contentFile string, address string) string {
//	conn := getTcpConnection(address)
//	defer conn.Close()
//	p := make([]byte, 2048)
//	lenData := getMockData(p, contentFile)
//	sendLen, sendErr := conn.Write(p[:lenData])
//	gomega.Expect(sendErr).ShouldNot(gomega.HaveOccurred())
//	gomega.Expect(sendLen).To(gomega.BeEquivalentTo(lenData))
//	time.Sleep(100 * time.Millisecond)
//	return conn.LocalAddr().String()
//}
//
//
//func TestSaveMsg(test *testing.T) {
//	gomega.RegisterTestingT(test)
//
//	fmt.Println("Trying to save a message on queue")
//
//	mu := sync.Mutex{}
//	q := getQueue(5)
//
//	mu.Lock()
//	count := 0
//	mu.Unlock()
//	q.OnPublishedEvent = func(message interface{}) {
//		mu.Lock()
//		defer mu.Unlock()
//		count++
//		fmt.Println("Publish OK")
//	}
//	mu.Lock()
//	countErr := 0
//	mu.Unlock()
//	q.OnNotPublishedEvent = func(message interface{}) {
//		mu.Lock()
//		defer mu.Unlock()
//		countErr++
//		fmt.Println("Publish Err")
//	}
//
//	err := q.StartConsume(func(queueName string, msg []byte) bool {
//		fmt.Println("New message")
//		return true
//	})
//	gomega.Expect(err).To(gomega.BeNil())
//
//	publishMockMessageOnQueue(q)
//
//	gomega.Eventually(func() int {
//		mu.Lock()
//		defer mu.Unlock()
//		return count
//	}).Should(gomega.BeEquivalentTo(1))
//	gomega.Eventually(func() int {
//		mu.Lock()
//		defer mu.Unlock()
//		return countErr
//	}).Should(gomega.BeEquivalentTo(0))
//
//	time.Sleep(1 * time.Second)
//	q.Close()
//	q.WaitQueue()
//}
//
//func TestGetMsg(test *testing.T) {
//	gomega.RegisterTestingT(test)
//
//	const QTDE_MSGS = 1000
//
//	fmt.Println("Trying to consume messages on queue")
//
//	mu := sync.Mutex{}
//
//	q := getQueue(500)
//	defer q.Close()
//
//	mu.Lock()
//	newMessage := 0
//	mu.Unlock()
//
//	mu.Lock()
//	receivedMsg := Message{}
//	mu.Unlock()
//
//	err := q.StartConsume(func(queueName string, msg []byte) bool {
//		mu.Lock()
//		defer mu.Unlock()
//		newMessage++
//		err := json.Unmarshal(msg, &receivedMsg)
//		if err != nil {
//			return false
//		}
//		mockMsg := getMockMsg()
//		compareMsg(receivedMsg, mockMsg)
//		return true
//	})
//	gomega.Expect(err).To(gomega.BeNil())
//
//	mu.Lock()
//	newMessage = 0
//	mu.Unlock()
//
//	for i := 0; i < QTDE_MSGS; i++ {
//		publishMockMessageOnQueue(q)
//	}
//
//	gomega.Eventually(func() int {
//		mu.Lock()
//		defer mu.Unlock()
//		return newMessage
//	}, 2, 1).Should(gomega.BeEquivalentTo(QTDE_MSGS))
//
//	q.Close()
//	q.WaitQueue()
//}
//
//func TestDelayedGetMsgAtOnce(test *testing.T) {
//	gomega.RegisterTestingT(test)
//
//	const QTDE_MSGS = 1000
//
//	fmt.Println("Trying to consume messages on queue (delayed process)")
//
//	mu := sync.Mutex{}
//
//	q := getQueue(QTDE_MSGS)
//	defer q.Close()
//
//	mu.Lock()
//	newMessage := 0
//	mu.Unlock()
//
//	err := q.StartConsume(func(queueName string, msg []byte) bool {
//		defer func() {
//			mu.Lock()
//			newMessage++
//			mu.Unlock()
//		}()
//		receivedMsg := Message{}
//		err := json.Unmarshal(msg, &receivedMsg)
//		if err != nil {
//			return false
//		}
//		mockMsg := getMockMsg()
//		compareMsg(receivedMsg, mockMsg)
//		time.Sleep(750 * time.Millisecond)
//		return true
//	})
//	gomega.Expect(err).To(gomega.BeNil())
//
//	mu.Lock()
//	newMessage = 0
//	mu.Unlock()
//
//	for i := 0; i < QTDE_MSGS; i++ {
//		publishMockMessageOnQueue(q)
//	}
//	gomega.Eventually(func() int {
//		mu.Lock()
//		defer mu.Unlock()
//		return newMessage
//	}, 1).Should(gomega.BeEquivalentTo(QTDE_MSGS))
//	time.Sleep(time.Second)
//	q.Close()
//	q.WaitQueue()
//}
//
//func TestDelayedGetMsg(test *testing.T) {
//	gomega.RegisterTestingT(test)
//
//	const QTDE_MSGS = 1000
//
//	fmt.Println("Trying to consume messages on queue (delayed process)")
//
//	mu := sync.Mutex{}
//
//	q := getQueue(QTDE_MSGS / 2)
//	defer q.Close()
//
//	mu.Lock()
//	newMessage := 0
//	mu.Unlock()
//
//	err := q.StartConsume(func(queueName string, msg []byte) bool {
//		defer func() {
//			mu.Lock()
//			newMessage++
//			mu.Unlock()
//		}()
//		receivedMsg := Message{}
//		err := json.Unmarshal(msg, &receivedMsg)
//		if err != nil {
//			return false
//		}
//		mockMsg := getMockMsg()
//		compareMsg(receivedMsg, mockMsg)
//		time.Sleep(750 * time.Millisecond)
//		return true
//	})
//	gomega.Expect(err).To(gomega.BeNil())
//
//	mu.Lock()
//	newMessage = 0
//	mu.Unlock()
//
//	for i := 0; i < QTDE_MSGS; i++ {
//		publishMockMessageOnQueue(q)
//	}
//
//	gomega.Eventually(func() int {
//		mu.Lock()
//		defer mu.Unlock()
//		return newMessage
//	}, 2).Should(gomega.BeEquivalentTo(QTDE_MSGS))
//	time.Sleep(time.Second)
//	q.Close()
//	q.WaitQueue()
//}
