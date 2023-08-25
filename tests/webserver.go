package tests

import (
	"encoding/json"
	"fmt"
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
	"io"
	"sync"
	"time"
)

func CreateEventTester() {
	ginkgo.It("Valid CloudEvent msg", func() {
		testSendValidCloudEventRequest()
	})
	ginkgo.It("InValid CloudEvent msg", func() {
		testSendInvalidCloudEventRequest()
	})
	ginkgo.It("Valid CloudEvent (dataOnly) msg", func() {
		testSendValidCloudEventDataOnlyRequest()
	})
	ginkgo.It("Invalid CloudEvent (dataOnly) msg", func() {
		testSendInvalidCloudEventDataOnlyRequest()
	})
	ginkgo.It("Valid json msg", func() {
		testSendValidJsonRequest()
	})
}

func testSendValidCloudEventRequest() {
	fmt.Println("Sending a valid CloudEvent data")
	wg := sync.WaitGroup{}
	mockReader, err := GetMockReader(GetMockEvent(time.Now().UTC(), "CE", "1"))
	gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
	h := make(map[string]string)
	h[CONTENT_TYPE] = CONTENT_TYPE_CE
	mu := sync.Mutex{}
	mu.Lock()
	count := 0
	mu.Unlock()
	q := InitQueue(TEST_PUBLISH_QUEUE, &count, func(queueName string, msg []byte, counter int) bool {
		var mock MockEvent
		err := json.Unmarshal(msg, &mock)
		gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
		fmt.Println(fmt.Sprintf("cnt=%d, %s", counter, msg))
		fmt.Println(mock)
		return true
	})
	defer q.Close()
	wg.Add(1)
	go func() {
		defer ginkgo.GinkgoRecover()
		defer wg.Done()
		resp, err := SendPostRequestWithHeaders(TEST_ENDPOINT, mockReader, h)
		bodyStr, _ := io.ReadAll(resp.Body)
		println(string(bodyStr))
		gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
		gomega.Expect(resp.StatusCode).To(gomega.Equal(201))
	}()
	wg.Wait()
	gomega.Eventually(func() int {
		mu.Lock()
		defer mu.Unlock()
		return count
	}).Should(gomega.BeEquivalentTo(1))
}

func testSendInvalidCloudEventRequest() {
	fmt.Println("Sending an invalid CloudEvent data")
	wg := sync.WaitGroup{}
	mockEvent := GetMockEvent(time.Now().UTC(), "CE", "1")
	mockEvent.PublishQueue = "dummy_queue"
	mockReader, err := GetMockReader(mockEvent)
	gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
	h := make(map[string]string)
	h[CONTENT_TYPE] = CONTENT_TYPE_CE
	mu := sync.Mutex{}
	mu.Lock()
	count := 0
	mu.Unlock()
	q := InitQueue(TEST_PUBLISH_QUEUE, &count, func(queueName string, msg []byte, counter int) bool {
		//fmt.Println(fmt.Sprintf("cnt=%d, %s", counter, msg))
		return true
	})
	defer q.Close()
	wg.Add(1)
	go func() {
		defer ginkgo.GinkgoRecover()
		defer wg.Done()
		resp, err := SendPostRequestWithHeaders(TEST_ENDPOINT, mockReader, h)
		gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
		gomega.Expect(resp.StatusCode).To(gomega.Equal(500))
	}()
	wg.Wait()
	gomega.Consistently(func() int {
		mu.Lock()
		defer mu.Unlock()
		return count
	}).Should(gomega.BeEquivalentTo(0))
}

func testSendValidJsonRequest() {
	fmt.Println("Sending a valid json data")
	wg := sync.WaitGroup{}
	mockReader, err := GetMockReader(MockData{Ref: "1", PublishDate: time.Now().Format(DATE_FORMAT)})
	gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
	h := make(map[string]string)
	h[CONTENT_TYPE] = CONTENT_TYPE_JSON
	h[PUBLISH_QUEUE] = TEST_PUBLISH_QUEUE
	h[PUBLISH_TYPE] = TEST_PUBLISH_TYPE
	h[PUBLISH_DATE] = time.Now().Format(DATE_FORMAT)
	mu := sync.Mutex{}
	mu.Lock()
	count := 0
	mu.Unlock()
	q := InitQueue(TEST_PUBLISH_QUEUE, &count, func(queueName string, msg []byte, counter int) bool {
		var mock MockData
		err := json.Unmarshal(msg, &mock)
		gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
		//fmt.Println(fmt.Sprintf("cnt=%d, %s", counter, msg))
		//fmt.Println(mock)
		return true
	})
	defer q.Close()
	wg.Add(1)
	go func() {
		defer ginkgo.GinkgoRecover()
		defer wg.Done()
		resp, err := SendPostRequestWithHeaders(TEST_ENDPOINT, mockReader, h)
		gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
		gomega.Expect(resp.StatusCode).To(gomega.Equal(201))
	}()
	wg.Wait()
	gomega.Eventually(func() int {
		mu.Lock()
		defer mu.Unlock()
		return count
	}).Should(gomega.BeEquivalentTo(1))
}

func testSendValidCloudEventDataOnlyRequest() {
	fmt.Println("Sending a valid CloudEvent (dataOnly) data")
	wg := sync.WaitGroup{}
	mockReader, err := GetMockReader(GetMockEvent(time.Now().UTC(), TEST_PUBLISH_TYPE, "2"))
	gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
	h := make(map[string]string)
	h[CONTENT_TYPE] = CONTENT_TYPE_CE
	mu := sync.Mutex{}
	mu.Lock()
	count := 0
	mu.Unlock()
	q := InitQueue(TEST_PUBLISH_QUEUE, &count, func(queueName string, msg []byte, counter int) bool {
		//fmt.Println(fmt.Sprintf("cnt=%d, %s", counter, msg))
		return true
	})
	defer q.Close()
	wg.Add(1)
	go func() {
		defer ginkgo.GinkgoRecover()
		defer wg.Done()
		resp, err := SendPostRequestWithHeaders(TEST_ENDPOINT, mockReader, h)
		gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
		gomega.Expect(resp.StatusCode).To(gomega.Equal(201))
	}()
	wg.Wait()
	gomega.Eventually(func() int {
		mu.Lock()
		defer mu.Unlock()
		return count
	}).Should(gomega.BeEquivalentTo(1))
}

func testSendInvalidCloudEventDataOnlyRequest() {
	fmt.Println("Sending an invalid CloudEvent (dataOnly) data")
	wg := sync.WaitGroup{}
	mockEvent := GetMockEvent(time.Now().UTC(), "CE", "1")
	mockEvent.PublishQueue = "dummy_queue"
	mockReader, err := GetMockReader(mockEvent)
	gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
	h := make(map[string]string)
	h[CONTENT_TYPE] = CONTENT_TYPE_CE
	mu := sync.Mutex{}
	mu.Lock()
	count := 0
	mu.Unlock()
	q := InitQueue(TEST_PUBLISH_QUEUE, &count, func(queueName string, msg []byte, counter int) bool {
		//fmt.Println(fmt.Sprintf("cnt=%d, %s", counter, msg))
		return true
	})
	defer q.Close()
	wg.Add(1)
	go func() {
		defer ginkgo.GinkgoRecover()
		defer wg.Done()
		resp, err := SendPostRequestWithHeaders(TEST_ENDPOINT, mockReader, h)
		gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
		gomega.Expect(resp.StatusCode).To(gomega.Equal(500))
	}()
	wg.Wait()
	gomega.Consistently(func() int {
		mu.Lock()
		defer mu.Unlock()
		return count
	}).Should(gomega.BeEquivalentTo(0))
}
