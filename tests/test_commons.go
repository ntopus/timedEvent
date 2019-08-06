package tests

import (
	"github.com/onsi/gomega"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

func NewPostRequestWithHeaders(url string, data url.Values, headers map[string]string) (resp *http.Response, err error) {
	req, err := http.NewRequest("POST", url, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	for key, value := range headers {
		req.Header.Set(key, value)
	}
	client := http.DefaultClient
	return client.Do(req)
}

func sendGetRequest(url string) (resp *http.Response, err error) {
	client := http.Client{}
	req, err := http.NewRequest(http.MethodGet, url, nil)
	gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
	//req.Header.Set("token", GetTestToken())
	return client.Do(req)
}

func getConfigPath() string {
	cwd, err := os.Getwd()
	gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
	return filepath.Join(cwd, "test", "config.json")
}

//func SaveConfigFile(config file_config.AppConfig) {
//	_, err := os.Stat(filepath.Dir(getConfigPath()))
//	if os.IsNotExist(err) {
//		err := os.Mkdir(filepath.Dir(getConfigPath()), os.ModePerm)
//		gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
//	}
//	bytes, err := json.Marshal(config)
//	gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
//	err = ioutil.WriteFile(getConfigPath(), bytes, os.ModePerm)
//	gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
//}
//
//func RunApp() *gexec.Session {
//	command := exec.Command("make", "build")
//	err := command.Run()
//	gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
//	err = os.Setenv("MONGODB_SERVER", "fleet.db.interno.ntopus.com.br")
//	gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
//	binPath := filepath.Join(getBinPath(), "fleet-management-api")
//	command = exec.Command(binPath, "-c="+getConfigPath())
//	session, err := gexec.Start(command, ginkgo.GinkgoWriter, ginkgo.GinkgoWriter)
//	gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
//	time.Sleep(400 * time.Millisecond)
//	fmt.Println("Application is running")
//	return session
//}
