package tests

import (
	"fmt"
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
	"strconv"
	"sync"
	"time"
)

const (
	DATE_FORMAT        = "2006-01-02 15:04:05Z"
	TEST_ENDPOINT      = "http://localhost:9010/v1/event"
	CONTENT_TYPE       = "Content-Type"
	CONTENT_TYPE_CE    = "application/cloudevents"
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

func CreateEventRequest() {
	ginkgo.It("Valid CloudEvent msg", func() {
		fmt.Println("Sending a valid CloudEvent data")
		const TESTE_QTD = 10
		wg := sync.WaitGroup{}

		for i := 0; i < TESTE_QTD; i++ {
			strIvalue := strconv.Itoa(i)
			mockReader, err := GetMockReader(getMockEvent(time.Now().UTC(), "CE", strIvalue))
			gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
			h := make(map[string]string)
			h[CONTENT_TYPE] = CONTENT_TYPE_CE
			h[PUBLISH_DATE] = time.Now().Add(time.Duration(i) * time.Millisecond).UTC().Format(DATE_FORMAT)
			h[PUBLISH_QUEUE] = TEST_PUBLISH_QUEUE
			wg.Add(1)
			go func() {
				defer ginkgo.GinkgoRecover()
				defer wg.Done()
				resp, err := SendPostRequestWithHeaders(TEST_ENDPOINT, mockReader, h)
				gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
				gomega.Expect(resp.StatusCode).To(gomega.Equal(201))
			}()
			wg.Wait()
		}
		mu := sync.Mutex{}
		q := GetQueue(TEST_PUBLISH_QUEUE, 200)
		mu.Lock()
		count := 0
		mu.Unlock()
		err := q.StartConsume(func(queueName string, msg []byte) bool {
			mu.Lock()
			defer mu.Unlock()
			count++
			fmt.Println(fmt.Sprintf("cnt=%d, %s", count, msg))
			return true
		})
		gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
		gomega.Eventually(func() int {
			mu.Lock()
			defer mu.Unlock()
			return count
		}).Should(gomega.BeEquivalentTo(TESTE_QTD))
	})
	ginkgo.It("Valid CloudEvent (dataOnly) msg", func() {
		fmt.Println("Sending a valid CloudEvent (dataOnly) data")
		const TESTE_QTD = 10
		wg := sync.WaitGroup{}

		for i := 0; i < TESTE_QTD; i++ {
			strIvalue := strconv.Itoa(i)
			mockReader, err := GetMockReader(getMockEvent(time.Now().UTC(), TEST_PUBLISH_TYPE, strIvalue))
			gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
			h := make(map[string]string)
			h[CONTENT_TYPE] = CONTENT_TYPE_CE
			h[PUBLISH_DATE] = time.Now().Add(time.Duration(i) * time.Millisecond).UTC().Format(DATE_FORMAT)
			h[PUBLISH_QUEUE] = TEST_PUBLISH_QUEUE
			h[PUBLISH_TYPE] = TEST_PUBLISH_TYPE
			wg.Add(1)
			go func() {
				defer ginkgo.GinkgoRecover()
				defer wg.Done()
				resp, err := SendPostRequestWithHeaders(TEST_ENDPOINT, mockReader, h)
				gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
				gomega.Expect(resp.StatusCode).To(gomega.Equal(201))
			}()
			wg.Wait()
		}
		mu := sync.Mutex{}
		q := GetQueue(TEST_PUBLISH_QUEUE, 200)
		mu.Lock()
		countDO := 0
		mu.Unlock()
		err := q.StartConsume(func(queueName string, msg []byte) bool {
			mu.Lock()
			defer mu.Unlock()
			countDO++
			fmt.Println(fmt.Sprintf("cnt=%d, %s", countDO, msg))
			return true
		})
		gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
		gomega.Eventually(func() int {
			mu.Lock()
			defer mu.Unlock()
			return countDO
		}).Should(gomega.BeEquivalentTo(TESTE_QTD))
	})
}

func getMockEvent(publihsDate time.Time, publishType string, ref string) interface{} {
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
