package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/ivanmeca/timedEvent/application/modules/config"
	"github.com/ivanmeca/timedQueue/queue"
	queue_repository "github.com/ivanmeca/timedQueue/queue-repository"
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

func GetQueue(queueName string, threadLimit int) *queue.Queue {
	qr, err := queue_repository.NewQueueRepository(queue_repository.NewQueueRepositoryParams("randomUser", "randomPass", "srvqueue.module.ntopus.com.br", 5672))
	params := queue.NewQueueParams(queueName)
	params.SetThreadLimit(threadLimit)
	q, err := qr.QueueDeclare(params, true)
	gomega.Expect(err).To(gomega.BeNil())
	return q
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
