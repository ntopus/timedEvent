package data_types

import (
	"github.com/cloudevents/sdk-go/pkg/cloudevents"
	"github.com/pborman/uuid"
	"time"
)

type ArangoCloudEvent struct {
	ArangoId  string `json:"_id,omitempty"`
	ArangoKey string `json:"_key,omitempty"`
	ArangoRev string `json:"_rev,omitempty"`
	CloudEvent
}

func NewArangoCloudEventV02(eventType string, data string, extensions map[string]interface{}) (*ArangoCloudEvent, error) {
	e := &ArangoCloudEvent{}
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
