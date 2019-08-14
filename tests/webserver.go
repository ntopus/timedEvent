package tests

import (
	"fmt"
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
	"sync"
	"time"
)

func CreateEventTester() {
	ginkgo.It("Valid CloudEvent msg", func() {
		fmt.Println("Sending a valid CloudEvent data")
		wg := sync.WaitGroup{}
		mockReader, err := GetMockReader(getMockEvent(time.Now().UTC(), "CE", "1"))
		gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
		h := make(map[string]string)
		h[CONTENT_TYPE] = CONTENT_TYPE_CE
		wg.Add(1)
		go func() {
			defer ginkgo.GinkgoRecover()
			defer wg.Done()
			resp, err := SendPostRequestWithHeaders(TEST_ENDPOINT, mockReader, h)
			gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
			gomega.Expect(resp.StatusCode).To(gomega.Equal(201))
		}()
		wg.Wait()
		mu := sync.Mutex{}
		mu.Lock()
		count := 0
		mu.Unlock()
		Consumer = func(queueName string, msg []byte) bool {
			mu.Lock()
			defer mu.Unlock()
			count++
			fmt.Println(fmt.Sprintf("cnt=%d, %s", count, msg))
			return true
		}
		gomega.Eventually(func() int {
			mu.Lock()
			defer mu.Unlock()
			return count
		}).Should(gomega.BeEquivalentTo(1))
	})
	ginkgo.It("InValid CloudEvent msg", func() {
		fmt.Println("Sending an invalid CloudEvent data")
		wg := sync.WaitGroup{}
		mockEvent := getMockEvent(time.Now().UTC(), "CE", "1")
		mockEvent.PublishQueue = "dummy_queue"
		mockReader, err := GetMockReader(mockEvent)
		gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
		h := make(map[string]string)
		h[CONTENT_TYPE] = CONTENT_TYPE_CE
		wg.Add(1)
		go func() {
			defer ginkgo.GinkgoRecover()
			defer wg.Done()
			resp, err := SendPostRequestWithHeaders(TEST_ENDPOINT, mockReader, h)
			gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
			gomega.Expect(resp.StatusCode).To(gomega.Equal(500))
		}()
		wg.Wait()
		mu := sync.Mutex{}
		mu.Lock()
		count := 0
		mu.Unlock()
		Consumer = func(queueName string, msg []byte) bool {
			mu.Lock()
			defer mu.Unlock()
			count++
			fmt.Println(fmt.Sprintf("cnt=%d, %s", count, msg))
			return true
		}
		gomega.Consistently(func() int {
			mu.Lock()
			defer mu.Unlock()
			return count
		}).Should(gomega.BeEquivalentTo(0))
	})
	ginkgo.It("Valid CloudEvent (dataOnly) msg", func() {
		fmt.Println("Sending a valid CloudEvent (dataOnly) data")
		const TESTE_QTD = 10
		wg := sync.WaitGroup{}
		mockReader, err := GetMockReader(getMockEvent(time.Now().UTC(), TEST_PUBLISH_TYPE, "2"))
		gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
		h := make(map[string]string)
		h[CONTENT_TYPE] = CONTENT_TYPE_CE
		wg.Add(1)
		go func() {
			defer ginkgo.GinkgoRecover()
			defer wg.Done()
			resp, err := SendPostRequestWithHeaders(TEST_ENDPOINT, mockReader, h)
			gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
			gomega.Expect(resp.StatusCode).To(gomega.Equal(201))
		}()
		wg.Wait()
		mu := sync.Mutex{}
		mu.Lock()
		countDO := 0
		mu.Unlock()
		Consumer = func(queueName string, msg []byte) bool {
			mu.Lock()
			defer mu.Unlock()
			countDO++
			fmt.Println(fmt.Sprintf("cntD=%d, %s", countDO, msg))
			return true
		}
		gomega.Eventually(func() int {
			mu.Lock()
			defer mu.Unlock()
			return countDO
		}).Should(gomega.BeEquivalentTo(1))
	})
	ginkgo.It("Invalid CloudEvent (dataOnly) msg", func() {
		fmt.Println("Sending an invalid CloudEvent (dataOnly) data")
		wg := sync.WaitGroup{}
		mockEvent := getMockEvent(time.Now().UTC(), "CE", "1")
		mockEvent.PublishQueue = "dummy_queue"
		mockReader, err := GetMockReader(mockEvent)
		gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
		h := make(map[string]string)
		h[CONTENT_TYPE] = CONTENT_TYPE_CE
		wg.Add(1)
		go func() {
			defer ginkgo.GinkgoRecover()
			defer wg.Done()
			resp, err := SendPostRequestWithHeaders(TEST_ENDPOINT, mockReader, h)
			gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
			gomega.Expect(resp.StatusCode).To(gomega.Equal(500))
		}()
		wg.Wait()
		mu := sync.Mutex{}
		mu.Lock()
		countDO := 0
		mu.Unlock()
		Consumer = func(queueName string, msg []byte) bool {
			mu.Lock()
			defer mu.Unlock()
			countDO++
			fmt.Println(fmt.Sprintf("cntD=%d, %s", countDO, msg))
			return true
		}
		gomega.Consistently(func() int {
			mu.Lock()
			defer mu.Unlock()
			return countDO
		}).Should(gomega.BeEquivalentTo(0))
	})
}
