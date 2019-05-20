package arangoDB

import (
	"context"
	"errors"
	"fmt"
	"github.com/arangodb/go-driver"
	"github.com/ivanmeca/timedEvent/application/modules/database"
	"github.com/ivanmeca/timedEvent/application/modules/database/data_types"
)

type Collection struct {
	db               driver.Database
	collection       string
	collectionDriver driver.Collection
}

func (coll *Collection) DeleteItem(keyList []string) (bool, error) {
	var oldDocs []data_types.ArangoCloudEvent
	ctx := driver.WithReturnOld(context.Background(), oldDocs)
	for _, key := range keyList {
		_, err := coll.collectionDriver.RemoveDocument(ctx, key)
		if err != nil {
			return false, err
		}
	}
	return true, nil
}

func (coll *Collection) Insert(item *data_types.ArangoCloudEvent) (bool, error) {
	_, err := coll.collectionDriver.CreateDocument(nil, item)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (coll *Collection) Update(patch map[string]interface{}, key string) (bool, error) {
	_, err := coll.collectionDriver.UpdateDocument(nil, key, patch)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (coll *Collection) Read(filters []database.AQLComparator) ([]data_types.ArangoCloudEvent, error) {
	var item data_types.ArangoCloudEvent
	var list []data_types.ArangoCloudEvent

	bindVars := map[string]interface{}{}
	query := fmt.Sprintf("FOR item IN %s ", coll.collection)
	glueStr := "FILTER"
	bindVarsNames := 0
	for _, filter := range filters {
		bindVars[string('A'+bindVarsNames)] = filter.Value
		query += fmt.Sprintf(" %s item.%s %s @%s", glueStr, filter.Field, filter.Comparator, string('A'+bindVarsNames))
		glueStr = "AND"
		bindVarsNames++
	}
	query += fmt.Sprintf(" SORT item.Context.time DESC RETURN item")
	cursor, err := coll.db.Query(nil, query, bindVars)
	if err != nil {
		return nil, errors.New("internal error: " + err.Error())
	}
	for cursor.HasMore() == true {
		_, err = cursor.ReadDocument(nil, &item)
		if err != nil {
			return nil, errors.New("internal error: " + err.Error())
		}
		list = append(list, item)
	}
	defer cursor.Close()
	return list, nil
}

func (coll *Collection) ReadItem(key string) (*data_types.ArangoCloudEvent, error) {
	var item data_types.ArangoCloudEvent
	_, err := coll.collectionDriver.ReadDocument(nil, key, &item)
	if err != nil {
		return nil, err
	}
	return &item, nil
}
