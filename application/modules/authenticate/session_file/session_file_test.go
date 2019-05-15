package session_file

import (
	"github.com/onsi/gomega"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestHttpAuthenticateCreateTokenFile(t *testing.T) {
	gomega.RegisterTestingT(t)
	os.RemoveAll("./tmp_test/")
	var tc ISessionFile
	tc, err := NewSessionFile("./tmp_test/")
	gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
	err = tc.CreateTokenFile("token")
	gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
	defer os.RemoveAll("./tmp_test/")
	_, err = os.Stat("./tmp_test/token")
	gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
}

func TestHttpAuthenticateCheckTokenFile(t *testing.T) {
	gomega.RegisterTestingT(t)
	var tc ISessionFile
	tc, err := NewSessionFile("./test/")
	gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
	ok, err := tc.CheckTokenFile("token")
	gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
	gomega.Expect(ok).To(gomega.BeTrue())
}

func TestHttpAuthenticateCheckUncreatedTokenFile(t *testing.T) {
	gomega.RegisterTestingT(t)
	var tc ISessionFile
	tc, err := NewSessionFile("./tmp_test/")
	gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
	ok, err := tc.CheckTokenFile("token")
	gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
	gomega.Expect(ok).To(gomega.BeFalse())
}

func TestCreateAuthWrongFilePath(t *testing.T) {
	gomega.RegisterTestingT(t)
	os.RemoveAll("./tmp_test/")
	var tc ISessionFile
	tc, err := NewSessionFile("")
	gomega.Expect(err).Should(gomega.HaveOccurred())
	gomega.Expect(tc).To(gomega.BeNil())
}

func TestHttpAuthenticateCreateWrongTokenFile(t *testing.T) {
	gomega.RegisterTestingT(t)
	os.RemoveAll("./tmp_test/")
	var tc ISessionFile
	tc, err := NewSessionFile("./tmp_test/")
	gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
	err = tc.CreateTokenFile("")
	defer os.RemoveAll("./tmp_test/")
	gomega.Expect(err).Should(gomega.HaveOccurred())
	fi, err := ioutil.ReadDir("./tmp_test/")
	gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
	gomega.Expect(len(fi)).To(gomega.Equal(0))
}

func TestHttpAuthenticateDeleteTokenFile(t *testing.T) {
	gomega.RegisterTestingT(t)
	os.RemoveAll("./tmp_test/")
	defer os.RemoveAll("./tmp_test/")
	err := os.MkdirAll("./tmp_test/", os.ModePerm)
	gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
	err = ioutil.WriteFile(filepath.Join("./tmp_test/", "token"), []byte(""), os.ModePerm)
	gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
	var tc ISessionFile
	tc, err = NewSessionFile("./tmp_test/")
	gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
	err = tc.DeleteTokenFile("token")
	gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
	_, err = os.Stat("./tmp_test/token")
	gomega.Expect(os.IsNotExist(err)).To(gomega.BeTrue())
}
