package file_config

import (
	. "github.com/onsi/gomega"
	"os"
	"testing"
)

func TestLoadConfigFromJsonFile(t *testing.T) {
	RegisterTestingT(t)
	config, err := LoadConfig("./test/config.json")
	Expect(err).Should(BeNil())
	_, err = os.Getwd()
	Expect(err).Should(BeNil())
	expectConfig := AppConfig{}
	Expect(*config).Should(BeEquivalentTo(expectConfig))
}
