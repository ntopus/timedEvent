package main

import (
	"context"
	"fmt"
	"github.com/ivanmeca/timedEvent/application"
	"github.com/ivanmeca/timedEvent/tests"
	"github.com/onsi/gomega"
	"sync"
	"testing"
	"time"
)

func BenchmarkWebServer(b *testing.B) {
	gomega.RegisterTestingT(b)
	tests.BuildApplication()
	tests.SaveConfigFile()
	ctx := context.Background()
	appMan := application.NewApplicationManager(tests.GetConfigPath())
	err := appMan.RunApplication(ctx)
	gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
	time.Sleep(time.Second)
	defer func() {
		fmt.Println("Killing application")
		ctx.Done()
	}()

	b.Run("Publishing multiple events updates in parallel", func(b *testing.B) {
		gomega.RegisterTestingT(b)
		b.ReportAllocs()
		wg := sync.WaitGroup{}
		for n := 0; n < b.N; n++ {
			for ref := 0; ref < 2000; ref++ {
				wg.Add(1)
				go func(ref int) {
					defer wg.Done()
					h := make(map[string]string)
					h[tests.CONTENT_TYPE] = tests.CONTENT_TYPE_CE
					mockReader, err := tests.GetMockReader(tests.GetMockEvent(time.Now().UTC().Add(50*time.Second), "CE", "teste"))
					gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
					resp, err := tests.SendPostRequestWithHeaders(tests.TEST_ENDPOINT, mockReader, h)
					gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
					gomega.Expect(resp.StatusCode).To(gomega.Equal(201))
				}(ref)
			}
			wg.Wait()
		}
	})

	b.Run("Publishing multiple events updates in parallel", func(b *testing.B) {
		gomega.RegisterTestingT(b)
		b.ReportAllocs()
		wg := sync.WaitGroup{}
		for n := 0; n < b.N; n++ {
			for ref := 0; ref < 2000; ref++ {
				wg.Add(1)
				go func(ref int) {
					defer wg.Done()
					h := make(map[string]string)
					h[tests.CONTENT_TYPE] = tests.CONTENT_TYPE_CE
					mockReader, err := tests.GetMockReader(tests.GetMockEvent(time.Now().UTC().Add(50*time.Second), "CE", fmt.Sprintf("%d", ref)))
					gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
					resp, err := tests.SendPostRequestWithHeaders(tests.TEST_ENDPOINT, mockReader, h)
					gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
					gomega.Expect(resp.StatusCode).To(gomega.Equal(201))
				}(ref)
			}
			wg.Wait()
		}
	})
}

//func testSendMultiplesValidCloudEventRequest() {
//	fmt.Println("Sending a valid CloudEvent data")
//	wg := sync.WaitGroup{}
//	mu := sync.Mutex{}
//	mu.Lock()
//	count := 0
//	mu.Unlock()
//	q := InitQueue(TEST_PUBLISH_QUEUE, &count, func(queueName string, msg []byte, counter int) bool {
//		var mock MockEvent
//		err := json.Unmarshal(msg, &mock)
//		gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
//		//fmt.Println(fmt.Sprintf("cnt=%d", counter))
//		//fmt.Println(mock)
//		return true
//	})
//	defer q.Close()
//	const TEST_QTDE = 10000
//	for i := 0; i < TEST_QTDE; i++ {
//		h := make(map[string]string)
//		h[CONTENT_TYPE] = CONTENT_TYPE_CE
//		wg.Add(1)
//		go func(ref string) {
//			defer ginkgo.GinkgoRecover()
//			defer wg.Done()
//			mockReader, err := GetMockReader(GetMockEvent(time.Now().UTC(), "CE", ref))
//			gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
//			resp, err := SendPostRequestWithHeaders(TEST_ENDPOINT, mockReader, h)
//			gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
//			gomega.Expect(resp.StatusCode).To(gomega.Equal(201))
//		}(fmt.Sprintf("%d", i))
//	}
//	wg.Wait()
//	gomega.Eventually(func() int {
//		mu.Lock()
//		defer mu.Unlock()
//		return count
//	}, 10).Should(gomega.BeEquivalentTo(TEST_QTDE))
//}
//
//func testSendMultiplesValidCloudEventUpdate() {
//	fmt.Println("Sending a valid CloudEvent data")
//	wg := sync.WaitGroup{}
//	mu := sync.Mutex{}
//	mu.Lock()
//	count := 0
//	mu.Unlock()
//	const TEST_QTDE = 1000
//	for i := 0; i < TEST_QTDE; i++ {
//		for ref := 0; ref < 10; ref++ {
//			h := make(map[string]string)
//			h[CONTENT_TYPE] = CONTENT_TYPE_CE
//			wg.Add(1)
//			go func(ref string) {
//				defer ginkgo.GinkgoRecover()
//				defer wg.Done()
//				mockReader, err := GetMockReader(GetMockEvent(time.Now().UTC().Add(50*time.Second), "CE", ref))
//				gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
//				resp, err := SendPostRequestWithHeaders(TEST_ENDPOINT, mockReader, h)
//				gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
//				gomega.Expect(resp.StatusCode).To(gomega.Equal(201))
//			}(fmt.Sprintf("%d", ref))
//		}
//	}
//	wg.Wait()
//	gomega.Eventually(func() int {
//		mu.Lock()
//		defer mu.Unlock()
//		return count
//	}, 10).Should(gomega.BeEquivalentTo(TEST_QTDE))
//}
