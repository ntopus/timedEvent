package tests

import (
	"github.com/onsi/ginkgo"
)

func CreateEventRequest() {
	ginkgo.It("Valid msg", func() {
		//for i := 1; i < 30; i++ {
		//	strIvalue := strconv.Itoa(i)
		//	fmt.Print("Trying to create a driver " + strIvalue)
		//	mockReader, err := getMockReader(getMockDriver(strIvalue))
		//	gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
		//	resp, err := sendPostRequest("http://localhost:8081/v1/driver?", mockReader)
		//	gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
		//	gomega.Expect(resp.StatusCode).To(gomega.Equal(201))
		//	respBody := parseResp(resp, &TestDriver)
		//	gomega.Expect(respBody).To(gomega.Not(gomega.BeNil()))
		//	compareDriver(getMockDriver(strIvalue), &TestDriver, false)
		//	fmt.Println(default_Tab_space, "SUCCESS")
		//}
	})
}
