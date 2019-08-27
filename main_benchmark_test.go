package main

import (
	"fmt"
	"github.com/ivanmeca/timedEvent/tests"
	"github.com/onsi/gomega"
	"testing"
	"time"
)

func BenchmarkApplication(b *testing.B) {
	gomega.RegisterTestingT(b)
	tests.BuildApplication()
	tests.SaveConfigFile()
	App = tests.RunApp()
	time.Sleep(time.Second)

	b.Run("Publishing multiple events updates in parallel", func(b *testing.B) {
		gomega.RegisterTestingT(b)
		b.ReportAllocs()
		for n := 0; n < b.N; n++ {
			for ref := 0; ref < 10; ref++ {
				h := make(map[string]string)
				h[tests.CONTENT_TYPE] = tests.CONTENT_TYPE_CE
				mockReader, err := tests.GetMockReader(tests.GetMockEvent(time.Now().UTC().Add(50*time.Second), "CE", fmt.Sprintf("%d", ref)))
				gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
				resp, err := tests.SendPostRequestWithHeaders(tests.TEST_ENDPOINT, mockReader, h)
				gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
				gomega.Expect(resp.StatusCode).To(gomega.Equal(201))
			}
		}
	})

	b.Run("Publishing multiple events in parallel", func(b *testing.B) {
		gomega.RegisterTestingT(b)
		b.ReportAllocs()
		for n := 0; n < b.N; n++ {
			h := make(map[string]string)
			h[tests.CONTENT_TYPE] = tests.CONTENT_TYPE_CE
			mockReader, err := tests.GetMockReader(tests.GetMockEvent(time.Now().UTC().Add(50*time.Second), "CE", fmt.Sprintf("%d", n)))
			gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
			resp, err := tests.SendPostRequestWithHeaders(tests.TEST_ENDPOINT, mockReader, h)
			gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
			gomega.Expect(resp.StatusCode).To(gomega.Equal(201))
		}
	})

	fmt.Println("Killing application")
	App.Kill()
}
