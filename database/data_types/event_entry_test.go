package data_types

import (
	"fmt"
	. "github.com/onsi/gomega"
	"testing"
)

func TestEventEntry(test *testing.T) {
	RegisterTestingT(test)
	fmt.Println("Trying to generate a event entry")
	event, err := NewCloudEventJsonV02(ApplicationJson, []byte(`"Teste"`), nil)
	Expect(err).ShouldNot(HaveOccurred())
	fmt.Println(event.String())
}

func TestEventEntryWithExtensions(test *testing.T) {
	RegisterTestingT(test)
	fmt.Println("Trying to generate a event entry with extensions")
	event, err := NewCloudEventJsonV02(ApplicationJson, []byte(`"Teste"`), map[string]interface{}{"DestPath": "teste"})
	Expect(err).ShouldNot(HaveOccurred())
	fmt.Println(event.String())
}
