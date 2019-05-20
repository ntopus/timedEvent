package data_types

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/cloudevents/sdk-go/pkg/cloudevents"
	"github.com/pborman/uuid"
	"strings"
	"time"
)

// Event represents the canonical representation of a CloudEvent.
type CloudEvent struct {
	Context     cloudevents.EventContextV02 `json:"context"`
	Data        interface{}                 `json:"data"`
	DataEncoded bool                        `json:"dataencoded"`
}

// New returns a new Event, an optional version can be passed to change the
// default spec version from 0.2 to the provided version.
func NewCloudEventV02(eventType string, data string, extensions map[string]interface{}) (*CloudEvent, error) {
	e := &CloudEvent{}
	err := e.Context.SetSpecVersion(cloudevents.CloudEventsVersionV02)
	if err != nil {
		return nil, err
	}
	err = e.Context.SetID(uuid.NewUUID().String())
	if err != nil {
		return nil, err
	}
	err = e.Context.SetType(eventType)
	if err != nil {
		return nil, err
	}
	err = e.Context.SetTime(time.Now())
	if err != nil {
		return nil, err
	}
	for key, value := range extensions {
		err = e.Context.SetExtension(key, value)
		if err != nil {
			return nil, err
		}
	}
	e.Data = data
	return e, nil
}

// ExtensionAs returns Context.ExtensionAs(name, obj)
func (e *CloudEvent) ExtensionAs(name string, obj interface{}) error {
	return e.Context.ExtensionAs(name, obj)
}

// Validate performs a spec based validation on this event.
// Validation is dependent on the spec version specified in the event context.
func (e *CloudEvent) Validate() error {
	if err := e.Context.Validate(); err != nil {
		return err
	}
	return nil
}

// String returns a pretty-printed representation of the Event.
func (e *CloudEvent) String() string {
	b := strings.Builder{}

	b.WriteString("Validation: ")

	valid := e.Validate()
	if valid == nil {
		b.WriteString("valid\n")
	} else {
		b.WriteString("invalid\n")
	}
	if valid != nil {
		b.WriteString(fmt.Sprintf("Validation Error: \n%s\n", valid.Error()))
	}

	b.WriteString(e.Context.String())

	if e.Data != "" {
		b.WriteString("Data,\n  ")
		if strings.HasPrefix(e.Context.GetDataContentType(), "application/json") {
			var prettyJSON bytes.Buffer
			data := e.Data.([]byte)
			err := json.Indent(&prettyJSON, data, "  ", "  ")
			if err != nil {
				b.Write([]byte(fmt.Sprintf("%s", e.Data)))
			} else {
				b.Write(prettyJSON.Bytes())
			}
		} else {
			b.Write([]byte(fmt.Sprintf("%s", e.Data)))
		}
		b.WriteString("\n")
	}
	return b.String()
}
