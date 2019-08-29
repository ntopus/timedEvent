package tests

import (
	"encoding/json"
	"fmt"
	"github.com/ivanmeca/timedEvent/application/modules/database/data_types"
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
	"sync"
	"time"
)

func SchedulerTester() {
	ginkgo.It("Test expired event", func() {
		testSendValidCloudEventRequestAndCheckDbContent()
	})
}

func testSendValidCloudEventRequestAndCheckDbContent() {
	fmt.Println("Sending a valid cloudEvent data")
	wg := sync.WaitGroup{}
	mu := sync.Mutex{}
	mu.Lock()
	count := 0
	mu.Unlock()
	q := InitQueue(TEST_PUBLISH_QUEUE, &count, func(queueName string, msg []byte, counter int) bool {
		defer ginkgo.GinkgoRecover()
		var mock MockData
		err := json.Unmarshal(msg, &mock)
		gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
		publishedDate, err := time.Parse("2006-01-02 15:04:05Z", mock.PublishDate)
		timeDiff := time.Now().UTC().Sub(publishedDate)
		gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
		fmt.Println(fmt.Sprintf("ref=%s,cnt=%d\tactualTime:%s\teventTime:%s\ttimeDiff: %v", mock.Ref, counter, time.Now().UTC().Format("15:04:05Z"), publishedDate.Format("15:04:05Z"), timeDiff))
		gomega.Expect(timeDiff).To(gomega.BeNumerically(">", 0))
		gomega.Expect(timeDiff).To(gomega.BeNumerically("<", 500*time.Millisecond))
		return true
	})
	defer q.Close()
	const TEST_QTDE = 10
	for i := 0; i < TEST_QTDE; i++ {
		h := make(map[string]string)
		h[CONTENT_TYPE] = CONTENT_TYPE_CE
		wg.Add(1)
		go func(ref int) {
			defer ginkgo.GinkgoRecover()
			defer wg.Done()
			delayToPublish := ref
			horaAtual := time.Now().UTC()
			mockReader, err := GetMockReader(GetMockEvent(horaAtual.Add(time.Duration(delayToPublish)*time.Second), data_types.DataOnly, fmt.Sprintf("%d", ref)))
			gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
			resp, err := SendPostRequestWithHeaders(TEST_ENDPOINT, mockReader, h)
			var MockEvent MockEvent
			ParseResp(resp, &MockEvent)
			gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
			gomega.Expect(resp.StatusCode).To(gomega.Equal(201))
			publishdDate, err := time.Parse("2006-01-02 15:04:05Z", MockEvent.PublishDate)
			gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
			fmt.Println(fmt.Sprintf("Actual time:%s\tEvent time:%s", time.Now().UTC().Format("15:04:05Z"), publishdDate.Format("15:04:05Z")))
		}(i)
	}
	wg.Wait()
	gomega.Eventually(func() int {
		mu.Lock()
		defer mu.Unlock()
		return count
	}, 10).Should(gomega.BeEquivalentTo(TEST_QTDE))
}
