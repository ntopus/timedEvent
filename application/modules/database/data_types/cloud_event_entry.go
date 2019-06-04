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
	cloudevents.EventContextV02
	PublishDate  time.Time   `json:"publishdate"`
	PublishQueue string      `json:"publishqueue"`
	Data         interface{} `json:"data"`
	DataEncoded  bool        `json:"dataencoded"`
}

// New returns a new Event, an optional version can be passed to change the
// default spec version from 0.2 to the provided version.
func NewCloudEventV02(eventType string, data interface{}, extensions map[string]interface{}) (*CloudEvent, error) {
	e := &CloudEvent{}
	err := e.SetSpecVersion(cloudevents.CloudEventsVersionV02)
	if err != nil {
		return nil, err
	}
	err = e.SetID(uuid.NewUUID().String())
	if err != nil {
		return nil, err
	}
	err = e.SetType(eventType)
	if err != nil {
		return nil, err
	}
	err = e.SetTime(time.Now())
	if err != nil {
		return nil, err
	}
	for key, value := range extensions {
		err = e.SetExtension(key, value)
		if err != nil {
			return nil, err
		}
	}
	e.Data = data
	return e, nil
}

// ExtensionAs returns Context.ExtensionAs(name, obj)
func (e *CloudEvent) ExtensionAs(name string, obj interface{}) error {
	return e.EventContextV02.ExtensionAs(name, obj)
}

// Validate performs a spec based validation on this event.
// Validation is dependent on the spec version specified in the event context.
func (e *CloudEvent) Validate() error {
	if err := e.EventContextV02.Validate(); err != nil {
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

	b.WriteString(e.EventContextV02.String())

	if e.Data != "" {
		b.WriteString("Data,\n  ")
		if strings.HasPrefix(e.GetDataContentType(), "application/json") {
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
