package data_types

import (
	"github.com/cloudevents/sdk-go/pkg/cloudevents"
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

type EventEntry struct {
	cloudevents.Event
}

func NewCloudEventJsonV02(eventType string, data string) *EventEntry {
	now := types.ParseTimestamp(time.Now().String())
	e := EventEntry{}
	e.Context = cloudevents.EventContextV02{
		ID:          uuid.NewUUID().String(),
		Type:        eventType,
		Time:        now,
		SpecVersion: cloudevents.CloudEventsVersionV02,
	}.AsV02()
	return &e
}

func (e *EventEntry) GetData() interface{} {
	return e.Data
}

func (e *EventEntry) GetType() string {
	return e.Context.GetType()
}

func (e *EventEntry) GetContentType() string {
	return e.Context.GetDataContentType()
}

func (e *EventEntry) GetSpecVersion() string {
	return e.Context.GetSpecVersion()
}

func (e *EventEntry) UnmarshalJSON(b []byte) error {
	return e.UnmarshalJSON(b)
}

func (e *EventEntry) MarshalJSON() ([]byte, error) {
	return e.MarshalJSON()
}
