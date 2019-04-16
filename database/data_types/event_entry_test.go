package data_types

import (
	"encoding/json"
	"fmt"
	"github.com/cloudevents/sdk-go"
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

func TestEventMarshal(test *testing.T) {
	RegisterTestingT(test)
	fmt.Println("Trying to marshalJSON a event entry")
	event, err := NewCloudEventJsonV02(ApplicationJson, []byte(`"Teste"`), map[string]interface{}{"DestPath": "teste"})
	data, err := json.Marshal(event)
	Expect(err).ShouldNot(HaveOccurred())
	fmt.Println(string(data))
	Expect(err).ShouldNot(HaveOccurred())
}

func TestNativeEventUnMarshal(test *testing.T) {
	RegisterTestingT(test)
	fmt.Println("Trying to UnmarshalJSON a event entry")
	data := `{"Context":{"specversion":"0.2","type":"application/json","source":"http://localhost:8080/","id":"97a58e52-5faa-11e9-a9ab-54bf64f7912d","time":"2019-04-15T18:16:30.479499099Z","-":{"DestPath":"teste"}},"Data":"IlRlc3RlIg=="}`
	var event cloudevents.Event
	err := json.Unmarshal([]byte(data), &event)
	Expect(err).ShouldNot(HaveOccurred())
	fmt.Println(event.String())
}

func TestEventUnMarshal(test *testing.T) {
	RegisterTestingT(test)
	fmt.Println("Trying to UnmarshalJSON a event entry")
	data := `{"Context":{"specversion":"0.2","type":"application/json","source":"http://localhost:8080/","id":"97a58e52-5faa-11e9-a9ab-54bf64f7912d","time":"2019-04-15T18:16:30.479499099Z","-":{"DestPath":"teste"}},"Data":"IlRlc3RlIg=="}`
	var event EventEntry
	err := json.Unmarshal([]byte(data), &event)
	Expect(err).ShouldNot(HaveOccurred())
	fmt.Println(event)
	fmt.Println(event.String())
}
