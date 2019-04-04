package data_types

import (
	"github.com/pborman/uuid"
	"time"
)

const (
	ApplicationJson = "application/json"
	DestinationPath = "destPath"
	DestinationType = "destType"
	DestinationTime = "destTime"
)

type EventEntry struct {
	v02.Event
}

func NewCloudEventJsonV02(eventType string, data string) *EventEntry {
	now := time.Now()
	e := &v02.Event{
		ContentType: ApplicationJson,
		Type:        eventType,
		ID:          uuid.NewUUID().String(),
		Time:        &now,
		SpecVersion: cloudevents.Version02,
		Data:        data,
	}
	return &EventEntry{*e}
}

func (e *EventEntry) GetData() interface{} {
	return e.Data
}

func (e *EventEntry) GetType() string {
	return e.Type
}

func (e *EventEntry) GetContentType() string {
	return e.ContentType
}

func (e *EventEntry) GetSpecVersion() string {
	return e.SpecVersion
}

func (e *EventEntry) UnmarshalJSON(b []byte) error {
	return e.Event.UnmarshalJSON(b)
}

func (e *EventEntry) MarshalJSON() ([]byte, error) {
	return e.Event.MarshalJSON()
}
