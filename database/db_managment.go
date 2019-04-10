package database

import (
	"errors"
	"fmt"
	"github.com/arangodb/go-driver"
	"github.com/ivanmeca/timedEvent/database/data_types"
)

type Management interface {
	Insert(item *data_types.EventEntry) (bool, error)
	ReadItem(key string, item *data_types.EventEntry) (bool, error)
	Update(patch map[string]interface{}, key string) (bool, error)
	HealthCheck() (bool, error)
	ReadCollection(collection string, filters map[string]interface{}) ([]data_types.EventEntry, error)
}

type EventDB struct {
	db   driver.Database
	coll map[string]driver.Collection
}

const (
	eventsCollection = "timedEvents"
)

func NewDbManagement(config DatabaseConfigurationReader) (Management, error) {
	client, err := NewClientDB(config)
	if err != nil {
		return nil, err
	}
	collections := []string{eventsCollection}
	db, collMap, err := schema(*client, config.GetDbName(), collections)
	if err != nil {
		return nil, err
	}
	return &EventDB{
		db:   db,
		coll: collMap,
	}, nil
}

func (db *EventDB) Insert(item *data_types.EventEntry) (bool, error) {
	coll, ok := db.coll[eventsCollection]
	if !ok {
		return false, errors.New("database collection error")
	}
	_, err := coll.CreateDocument(nil, item)
	if err != nil {
		return false, err
	}
	return true, nil
}
func (db *EventDB) ReadItem(key string, item *data_types.EventEntry) (bool, error) {
	coll, ok := db.coll[eventsCollection]
	if !ok {
		return false, errors.New("database collection error")
	}
	_, err := coll.ReadDocument(nil, key, item)
	if err != nil {
		return false, err
	}
	return true, nil
}
func (db *EventDB) Update(patch map[string]interface{}, key string) (bool, error) {
	coll, ok := db.coll[eventsCollection]
	if !ok {
		return false, errors.New("database collection error")
	}
	_, err := coll.UpdateDocument(nil, key, patch)
	if err != nil {
		return false, err
	}
	return true, nil
}
func (db *EventDB) HealthCheck() (bool, error) {
	_, err := db.db.Info(nil)
	if err != nil {
		return false, err
	}
	return true, nil
}
func (db *EventDB) ReadCollection(collection string, filters map[string]interface{}) ([]data_types.EventEntry, error) {
	query := fmt.Sprintf("FOR item IN %s FILTER", collection)
	glueStr := ""
	for key, value := range filters {
		query += fmt.Sprintf(" %s item.%s == @%s", glueStr, key, value)
		glueStr = "AND"
	}
	query += fmt.Sprintf("SORT item.Context.Time DESC RETURN l")
	cursor, err := db.db.Query(nil, query, filters)
	if err != nil {
		return []data_types.EventEntry{}, errors.New("internal error")
	}
	var object []data_types.EventEntry
	for cursor.HasMore() == true {
		var item data_types.EventEntry
		cursor.ReadDocument(nil, &item)
		object = append(object, item)
	}
	defer cursor.Close()
	return object, nil
}

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

func checkDateLayout(value string) string {
	var layout string
	if value[len(value)-1] == 'Z' {
		layout = "2006-01-02 15:04:05Z"
	} else {
		layout = "2006-01-02 15:04:05"
	}
	return layout
}
