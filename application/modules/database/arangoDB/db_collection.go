package arangoDB

import (
	"context"
	"fmt"
	"github.com/arangodb/go-driver"
	"github.com/ivanmeca/timedEvent/application/modules/database"
	"github.com/ivanmeca/timedEvent/application/modules/database/data_types"
	"github.com/pkg/errors"
)

type Collection struct {
	db               driver.Database
	collection       string
	collectionDriver driver.Collection
}

func (coll *Collection) DefaultErrorHandler(err error) error {
	return errors.Wrap(err, "internal db error")
}

func (coll *Collection) DeleteItem(keyList []string) ([]data_types.ArangoCloudEvent, error) {
	for _, key := range keyList {
		_, err := coll.collectionDriver.RemoveDocument(nil, key)
		if err != nil {
			return nil, err
		}
	}
	return nil, nil
}

func (coll *Collection) Insert(item *data_types.ArangoCloudEvent) (*data_types.ArangoCloudEvent, error) {
	var newDoc data_types.ArangoCloudEvent
	ctx := driver.WithReturnNew(context.Background(), &newDoc)
	_, err := coll.collectionDriver.CreateDocument(ctx, item)
	if err != nil {
		return nil, coll.DefaultErrorHandler(err)
	}
	return &newDoc, nil
}

func (coll *Collection) Upsert(item *data_types.ArangoCloudEvent) (*data_types.ArangoCloudEvent, error) {
	var newDoc data_types.ArangoCloudEvent
	ctx := driver.WithReturnNew(context.Background(), &newDoc)
	bindVars := map[string]interface{}{}
	bindVars["item"] = item
	//query := fmt.Sprintf(`UPSERT {_key:'%s'} INSERT @item REPLACE @item in %s OPTIONS { exclusive: true } RETURN NEW `, item.ArangoKey, coll.collection)
	query := fmt.Sprintf(`INSERT @item INTO %s OPTIONS { overwrite: true, exclusive: true } RETURN NEW `, coll.collection)
	cursor, err := coll.db.Query(ctx, query, bindVars)
	if err != nil {
		return nil, coll.DefaultErrorHandler(err)
	}
	defer cursor.Close()
	for cursor.HasMore() == true {
		_, err = cursor.ReadDocument(nil, &newDoc)
		if err != nil {
			return nil, coll.DefaultErrorHandler(err)
		}
	}
	return &newDoc, nil
}

func (coll *Collection) Update(patch map[string]interface{}, key string) (*data_types.ArangoCloudEvent, error) {
	var newDoc data_types.ArangoCloudEvent
	ctx := driver.WithReturnNew(context.Background(), &newDoc)
	_, err := coll.collectionDriver.UpdateDocument(ctx, key, patch)
	if err != nil {
		return nil, coll.DefaultErrorHandler(err)
	}
	return &newDoc, nil
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
	query += fmt.Sprintf(" SORT item.time DESC RETURN item")
	cursor, err := coll.db.Query(nil, query, bindVars)
	defer cursor.Close()
	if err != nil {
		return nil, coll.DefaultErrorHandler(err)
	}
	for cursor.HasMore() == true {
		_, err = cursor.ReadDocument(nil, &item)
		if err != nil {
			return nil, coll.DefaultErrorHandler(err)
		}
		list = append(list, item)
	}
	return list, nil
}

func (coll *Collection) ReadItem(key string) (*data_types.ArangoCloudEvent, error) {
	var item data_types.ArangoCloudEvent
	_, err := coll.collectionDriver.ReadDocument(nil, key, &item)
	if err != nil {
		return nil, coll.DefaultErrorHandler(err)
	}
	return &item, nil
}
