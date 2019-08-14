package tests

import (
	"fmt"
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
	"strconv"
	"sync"
	"time"
)

type fnConsume func(queueName string, msg []byte) bool

var Consumer fnConsume

func CreateEventTester() {
	ginkgo.It("Valid CloudEvent msg", func() {
		fmt.Println("Sending a valid CloudEvent data")
		wg := sync.WaitGroup{}
		mockReader, err := GetMockReader(getMockEvent(time.Now().UTC(), "CE", "1"))
		gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
		h := make(map[string]string)
		h[CONTENT_TYPE] = CONTENT_TYPE_CE
		h[PUBLISH_DATE] = time.Now().Add(time.Duration(2) * time.Second).UTC().Format(DATE_FORMAT)
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
		}).Should(gomega.BeEquivalentTo(TESTE_QTD))
	})
}
