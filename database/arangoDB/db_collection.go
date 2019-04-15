package arangoDB

import (
	"errors"
	"fmt"
	"github.com/arangodb/go-driver"
	"github.com/ivanmeca/timedEvent/database/data_types"
)

type ArangoDbCollection struct {
	db               driver.Database
	collection       string
	collectionDriver driver.Collection
}

func (coll *ArangoDbCollection) DeleteItem(keyList []string) (bool, error) {
	for _, key := range keyList {
		_, err := coll.collectionDriver.RemoveDocument(nil, key)
		if err != nil {
			return false, err
		}
	}
	return true, nil
}

func (coll *ArangoDbCollection) Insert(item *data_types.EventEntry) (bool, error) {
	_, err := coll.collectionDriver.CreateDocument(nil, item)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (coll *ArangoDbCollection) Update(patch map[string]interface{}, key string) (bool, error) {
	_, err := coll.collectionDriver.UpdateDocument(nil, key, patch)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (coll *ArangoDbCollection) Read(filters map[string]interface{}) ([]data_types.EventEntry, error) {
	var item data_types.EventEntry
	var list []data_types.EventEntry

	query := fmt.Sprintf("FOR item IN %s ", coll.collection)
	glueStr := "FILTER"
	for key, value := range filters {
		query += fmt.Sprintf(" %s item.%s == @%s", glueStr, key, value)
		glueStr = "AND"
	}
	query += fmt.Sprintf(" SORT item.Context.Time DESC RETURN item")
	cursor, err := coll.db.Query(nil, query, filters)
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

func (coll *ArangoDbCollection) ReadItem(key string) (*data_types.EventEntry, error) {
	var item data_types.EventEntry
	_, err := coll.collectionDriver.ReadDocument(nil, key, &item)
	if err != nil {
		return nil, err
	}
	return &item, nil
}

//func (db *EventDB) ReadCollection(collection string, filters map[string]interface{}) ([]data_types.EventEntry, error) {
//	query := fmt.Sprintf("FOR item IN %s ", collection)
//	glueStr := "FILTER"
//	for key, value := range filters {
//		query += fmt.Sprintf(" %s item.%s == @%s", glueStr, key, value)
//		glueStr = "AND"
//	}
//	query += fmt.Sprintf(" SORT item.Context.Time DESC RETURN item")
//	cursor, err := db.db.Query(nil, query, filters)
//	if err != nil {
//		return []data_types.EventEntry{}, errors.New("internal error")
//	}
//	var object []data_types.EventEntry
//	for cursor.HasMore() == true {
//		var item data_types.EventEntry
//		cursor.ReadDocument(nil, &item)
//		object = append(object, item)
//	}
//	defer cursor.Close()
//	return object, nil
//}

//func (l *LocationDB) ReadLocationWithFilters(location structs.Location) (structs.Location, error) {
//	bindVars := map[string]interface{}{
//		"_key":                  location.Key,
//		"deviceId":              location.DeviceId,
//		"systemType":            location.SystemType,
//		"system":                location.System,
//		"source":                location.Source,
//		"creationDate":          location.CreationDate,
//		"locationDate":          location.LocationDate,
//		"coordinates":           location.Coordinates,
//		"velocity":              location.Velocity,
//		"precision":             location.Precision,
//		"userId":                location.UserId,
//		"userDescription":       location.UserDescription,
//		"contractId":            location.ContractId,
//		"deviceUniqueId":        location.DeviceUniqueId,
//		"userTypeId":            location.UserTypeId,
//		"userTypeDescription":   location.UserTypeDescription,
//		"contractCode":          location.ContractCode,
//		"enterpriseId":          location.EnterpriseId,
//		"enterpriseDescription": location.EnterpriseDescription,
//		"systemId":              location.SystemId,
//		"systemDescription":     location.SystemDescription,
//		"systemTypeId":          location.SystemTypeId,
//		"systemTypeDescription": location.SystemTypeDescription,
//		"timeLapse":             location.TimeLapse,
//		"precisionQualifier":    location.PrecisionQualifier,
//		"reason":                location.Reason,
//		"version":               location.Version,
//		"odometer":              location.Odometer,
//		"horimetre":             location.HourMeter,
//		"batteryLevel":          location.BatteryLevel,
//		"extraFields":           location.ExtraFields,
//		"integrator":            location.Integrator,
//		"systemIntegrator":      location.SystemIntegrator,
//	}
//	query := "FOR l IN locationCollection FILTER "
//	glueStr := ""
//	for key := range bindVars {
//		query += fmt.Sprintf("%sl.%s == @%s", glueStr, key, key)
//		glueStr = " AND "
//	}
//	query += " RETURN l"
//	cursor, err := l.db.Query(nil, query, bindVars)
//	if err != nil {
//		return structs.Location{}, errors.New("internal error")
//	}
//	object := structs.Location{}
//	_, err = cursor.ReadDocument(nil, &object)
//	if err != nil {
//		return structs.Location{}, err
//	}
//	defer cursor.Close()
//	return object, nil
//}
