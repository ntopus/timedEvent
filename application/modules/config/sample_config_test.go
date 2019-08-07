package config

import (
	"fmt"
	"github.com/onsi/gomega"
	"testing"
	"time"
)

func TestConfigSample(test *testing.T) {
	gomega.RegisterTestingT(test)
	fmt.Println("Trying to generate a config sample file")
	coll := getTestCollectionInstance("testeCollection")
	horaAtual := time.Now().AddDate(0, 0, 3)
	list, err := coll.Read([]database.AQLComparator{{Field: "Context.time", Comparator: ">=", Value: horaAtual}})
	gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
	fmt.Println(list)
}
