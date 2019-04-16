package data_types

import (
	"encoding/json"
	"errors"
	"github.com/cloudevents/sdk-go/pkg/cloudevents/types"
	"github.com/pborman/uuid"
	"time"
)

const (
	ApplicationJson = "application/json"
	DestinationPath = "destPath"
	DestinationType = "destType"
	DestinationTime = "destTime"
)

const (
	// CloudEventsVersionV02 represents the version 0.2 of the CloudEvents spec.
	CloudEventsVersionV02 = "0.2"
)

type EventContext struct {
	// The version of the CloudEvents specification used by the event.
	SpecVersion string `json:"specversion"`
	// The type of the occurrence which has happened.
	Type string `json:"type"`
	// A URI describing the event producer.
	Source types.URLRef `json:"source"`
	// ID of the event; must be non-empty and unique within the scope of the producer.
	ID string `json:"id"`
	// Timestamp when the event happened.
	Time *time.Time `json:"time,omitempty"`
	// A link to the schema that the `data` attribute adheres to.
	SchemaURL *types.URLRef `json:"schemaurl,omitempty"`
	// A MIME (RFC2046) string describing the media type of `data`.
	ContentType *string `json:"contenttype,omitempty"`
	// Additional extension metadata beyond the base spec.
	Extensions map[string]interface{} `json:"-,omitempty"`
}

type EventEntry struct {
	Context     EventContext
	Data        interface{}
	DataEncoded bool
}

func NewCloudEventJsonV02(eventType string, data []byte, extensions map[string]interface{}) (*EventEntry, error) {

	var e EventEntry

	e.Context.SpecVersion = CloudEventsVersionV02
	e.Context.ID = uuid.NewUUID().String()
	e.Context.Type = eventType
	hora := time.Now()
	e.Context.Time = &hora
	e.Context.Extensions = map[string]interface{}{}
	for key, value := range extensions {
		e.Context.Extensions[key] = value
	}

	validJson := json.Valid(data)
	if !validJson {
		return nil, errors.New("invalid input data json")
	}
	e.Data = data
	return &e, nil
}

func (e *EventEntry) GetData() interface{} {
	return e.Data
}

func (e *EventEntry) GetType() string {
	return e.Context.Type
}

func (e *EventEntry) GetContentType() string {
	return *e.Context.ContentType
}

func (e *EventEntry) GetSpecVersion() string {
	return e.Context.SpecVersion
}

func (e *EventEntry) String() string {
	return ""
}

//func (e *EventEntry) UnmarshalJSON(b []byte) error {
//	var Event cloudevents.Event
//	err := json.Unmarshal(b, &Event)
//	if err != nil {
//		return err
//	}
//	e.ev = &Event
//	return nil
//}
//
//func (e *EventEntry) MarshalJSON() ([]byte, error) {
//	return json.Marshal(e.ev)
//}
