package data_types

import (
	"encoding/json"
	"errors"
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

func NewCloudEventJsonV02(eventType string, data []byte, extensions map[string]interface{}) (*EventEntry, error) {
	now := types.ParseTimestamp(time.Now().UTC().Format(time.RFC3339Nano))
	e := EventEntry{}
	e.Context = cloudevents.EventContextV02{
		ID:          uuid.NewUUID().String(),
		Type:        eventType,
		Time:        now,
		Source:      *types.ParseURLRef("http://localhost:8080/"),
		SpecVersion: cloudevents.CloudEventsVersionV02,
		Extensions:  extensions,
	}.AsV02()

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
	return e.Context.GetType()
}

func (e *EventEntry) GetContentType() string {
	return e.Context.GetDataContentType()
}

func (e *EventEntry) GetSpecVersion() string {
	return e.Context.GetSpecVersion()
}
