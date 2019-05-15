package data_types

import (
	"fmt"
	. "github.com/onsi/gomega"
	"testing"
)

func TestCloudEventEntry(test *testing.T) {
	RegisterTestingT(test)
	fmt.Println("Trying to generate a event entry")
	event, err := NewCloudEventV02(ApplicationJson, "Teste", nil)
	event.Context.SetSource("timedEvent")
	Expect(err).ShouldNot(HaveOccurred())
	fmt.Println(event.String())
}
