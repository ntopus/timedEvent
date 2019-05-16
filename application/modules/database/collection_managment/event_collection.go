package collection_managment

import (
	"github.com/ivanmeca/timedEvent/application/modules/database"
	"github.com/ivanmeca/timedEvent/application/modules/database/data_types"
)

type EventCollection struct {
}

func (e *EventCollection) Insert(item *data_types.CloudEvent) (bool, error) {
	return true, nil
}
func (e *EventCollection) DeleteItem(keyList []string) (bool, error) {
	return true, nil
}
func (e *EventCollection) Update(patch map[string]interface{}, key string) (bool, error) {
	return true, nil
}
func (e *EventCollection) Read(filters []database.AQLComparator) ([]data_types.CloudEvent, error) {
	return nil, nil
}
func (e *EventCollection) ReadItem(key string) (*data_types.CloudEvent, error) {
	return nil, nil
}
