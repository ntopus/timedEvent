package collection_managment

import (
	"github.com/ivanmeca/timedEvent/application/modules/database"
	"github.com/ivanmeca/timedEvent/application/modules/database/data_types"
)

const EventCollectionName = "TesteColl"

func NewEventCollection() *EventCollection {
	return &EventCollection{}
}

type EventCollection struct {
}

func (e *EventCollection) Insert(item *data_types.CloudEvent) (bool, error) {
	coll, err := GetDBSession().GetCollection(EventCollectionName)
	if err != nil {
		return false, err
	}
	event := data_types.ArangoCloudEvent{}
	event.ArangoId = item.Context.GetID()
	event.Context = item.Context
	event.Data = item.Data
	event.DataEncoded = item.DataEncoded
	return coll.Insert(&event)
}
func (e *EventCollection) DeleteItem(keyList []string) (bool, error) {
	coll, err := GetDBSession().GetCollection(EventCollectionName)
	if err != nil {
		return false, err
	}
	data, err := e.Read([]database.AQLComparator{{Field: "Context.id", Comparator: "==", Value: key}})
	if err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, nil
	}
	return &data[0], nil

	return coll.DeleteItem(keyList)
}
func (e *EventCollection) Update(patch map[string]interface{}, key string) (bool, error) {

	return true, nil
}
func (e *EventCollection) Read(filters []database.AQLComparator) ([]data_types.CloudEvent, error) {
	coll, err := GetDBSession().GetCollection(EventCollectionName)
	if err != nil {
		return nil, err
	}
	collectionData, err := coll.Read(filters)
	if err != nil {
		return nil, err
	}
	var returnData []data_types.CloudEvent
	for _, value := range collectionData {
		returnData = append(returnData, value.CloudEvent)
	}
	return returnData, nil
}
func (e *EventCollection) ReadItem(key string) (*data_types.CloudEvent, error) {
	data, err := e.Read([]database.AQLComparator{{Field: "Context.id", Comparator: "==", Value: key}})
	if err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, nil
	}
	return &data[0], nil
}
