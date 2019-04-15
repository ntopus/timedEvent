package data_types

import (
	"encoding/json"
	"errors"
	"github.com/cloudevents/sdk-go"
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
	ev *cloudevents.Event
}

func NewCloudEventJsonV02(eventType string, data []byte, extensions map[string]interface{}) (*EventEntry, error) {

	e := cloudevents.NewEvent(cloudevents.VersionV02)
	e.SetID(uuid.NewUUID().String())
	e.SetType(eventType)
	e.SetTime(time.Now())

	for key, value := range extensions {
		e.SetExtension(key, value)
	}

	validJson := json.Valid(data)
	if !validJson {
		return nil, errors.New("invalid input data json")
	}
	e.Data = data

	event := EventEntry{
		ev: &e,
	}

	return &event, nil
}

func (e *EventEntry) GetData() interface{} {
	return e.ev.Data
}

func (e *EventEntry) GetType() string {
	return e.ev.Context.GetType()
}

func (e *EventEntry) GetContentType() string {
	return e.ev.Context.GetDataContentType()
}

func (e *EventEntry) GetSpecVersion() string {
	return e.ev.Context.GetSpecVersion()
}

func (e *EventEntry) String() string {
	return e.ev.String()
}

func (e *EventEntry) UnmarshalJSON(b []byte) error {
	var Event cloudevents.Event
	err := json.Unmarshal(b, &Event)
	if err != nil {
		return err
	}
	e.ev = &Event
	return nil
}

func (e *EventEntry) MarshalJSON() ([]byte, error) {
	return json.Marshal(e.ev)
}
