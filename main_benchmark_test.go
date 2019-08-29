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
	appMan := application.NewApplicationManager("./config.json")
	err := appMan.RunApplication(ctx)
	gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
	time.Sleep(time.Second)
	defer func() {
		fmt.Println("Killing application")
		ctx.Done()
	}()

	//b.Run("Publishing multiple events updates in parallel", func(b *testing.B) {
	//	gomega.RegisterTestingT(b)
	//	b.ReportAllocs()
	//	for n := 0; n < b.N; n++ {
	//		h := make(map[string]string)
	//		h[tests.CONTENT_TYPE] = tests.CONTENT_TYPE_CE
	//		mockReader, err := tests.GetMockReader(tests.GetMockEvent(time.Now().UTC().Add(50*time.Second), "CE", "teste"))
	//		gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
	//		resp, err := tests.SendPostRequestWithHeaders(tests.TEST_ENDPOINT, mockReader, h)
	//		gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
	//		gomega.Expect(resp.StatusCode).To(gomega.Equal(201))
	//	}
	//})

	//b.Run("Publishing multiple events updates in parallel", func(b *testing.B) {
	//	gomega.RegisterTestingT(b)
	//	b.ReportAllocs()
	//	for n := 0; n < b.N; n++ {
	//		for ref := 0; ref < 10; ref++ {
	//			h := make(map[string]string)
	//			h[tests.CONTENT_TYPE] = tests.CONTENT_TYPE_CE
	//			mockReader, err := tests.GetMockReader(tests.GetMockEvent(time.Now().UTC().Add(50*time.Second), "CE", fmt.Sprintf("%d", ref)))
	//			gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
	//			resp, err := tests.SendPostRequestWithHeaders(tests.TEST_ENDPOINT, mockReader, h)
	//			gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
	//			gomega.Expect(resp.StatusCode).To(gomega.Equal(201))
	//		}
	//	}
	//})

	b.Run("Publishing multiple events in parallel", func(b *testing.B) {
		wg := sync.WaitGroup{}
		gomega.RegisterTestingT(b)
		b.ReportAllocs()
		for n := 0; n < b.N; n++ {
			wg.Add(1)
			go func(ref int) {
				h := make(map[string]string)
				h[tests.CONTENT_TYPE] = tests.CONTENT_TYPE_CE
				mockReader, err := tests.GetMockReader(tests.GetMockEvent(time.Now().UTC().Add(time.Second), "CE", fmt.Sprintf("%d", ref)))
				gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
				resp, err := tests.SendPostRequestWithHeaders(tests.TEST_ENDPOINT, mockReader, h)
				gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
				gomega.Expect(resp.StatusCode).To(gomega.Equal(201))
			}(n)
		}
	})

}
