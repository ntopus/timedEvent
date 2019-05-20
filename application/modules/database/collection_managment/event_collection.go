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
	return coll.Insert(item)
}
func (e *EventCollection) DeleteItem(keyList []string) (bool, error) {
	coll, err := GetDBSession().GetCollection(EventCollectionName)
	if err != nil {
		return false, err
	}
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
	return coll.Read(filters)
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
