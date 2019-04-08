package data_types

import (
	"fmt"
	. "github.com/onsi/gomega"
	"testing"
)

func TestEvbentEntry(test *testing.T) {
	RegisterTestingT(test)

	fmt.Println("Trying to generate a event entry")

	event, err := NewCloudEventJsonV02(ApplicationJson, []byte(`{"name":"Teste"}`))
	Expect(err).ShouldNot(HaveOccurred())
	fmt.Println(event.String())
}
