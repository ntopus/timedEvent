package data_types

import (
	"github.com/cloudevents/sdk-go"
	"github.com/cloudevents/sdk-go/v02"
	"github.com/pborman/uuid"
	"time"
)

const (
	ApplicationJson = "application/json"
	DestinationPath = "destPath"
	DestinationType = "destType"
	DestinationTime = "destTime"
)

type ServiceEntry struct {
	v02.Event
}

func NewCloudEventJsonV02(eventType string, data string) *ServiceEntry {
	now := time.Now()
	e := &v02.Event{
		ContentType: ApplicationJson,
		Type:        eventType,
		ID:          uuid.NewUUID().String(),
		Time:        &now,
		SpecVersion: cloudevents.Version02,
		Data:        data,
	}
	return &ServiceEntry{*e}
}

func (e *ServiceEntry) GetData() interface{} {
	return e.Data
}

func (e *ServiceEntry) GetType() string {
	return e.Type
}

func (e *ServiceEntry) GetContentType() string {
	return e.ContentType
}

func (e *ServiceEntry) GetSpecVersion() string {
	return e.SpecVersion
}

func (e *ServiceEntry) UnmarshalJSON(b []byte) error {
	return e.Event.UnmarshalJSON(b)
}

func (e *ServiceEntry) MarshalJSON() ([]byte, error) {
	return e.Event.MarshalJSON()
}
