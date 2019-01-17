package database

import (
	"devgit.kf.com.br/border/location-api/config"
	"devgit.kf.com.br/border/location-api/routes/helpers/constants"
	"devgit.kf.com.br/border/location-api/structs"
	"errors"
	"fmt"
	"github.com/arangodb/go-driver"
	"strconv"
	"time"
)

type Management interface {
	InsertLocation(location *structs.Location) (bool, error)
	ReadLocation(key string, location *structs.Location) (bool, error)
	UpdateLocation(patch map[string]interface{}, key string) (bool, error)
	HealthCheck() (bool, error)
	ReadLocationWithFilters(location structs.Location) (structs.Location, error)
	ReadLocations(filters map[string]string) ([]structs.Location, error)
}

type LocationDB struct {
	db   driver.Database
	coll map[string]driver.Collection
}

const (
	dbName             = "location"
	locationCollection = "locationCollection"
	statusCollection   = "statusCollection"
)

func NewDbManagement(config *config.AppConfig) (Management, error) {
	client, err := NewClientDB(config)
	if err != nil {
		return nil, err
	}
	collections := []string{locationCollection, statusCollection}
	db, collMap, err := schema(*client, collections)
	if err != nil {
		return nil, err
	}
	return &LocationDB{
		db:   db,
		coll: collMap,
	}, nil
}

func (l *LocationDB) ReadLocation(key string, location *structs.Location) (bool, error) {
	coll, ok := l.coll["locationCollection"]
	if !ok {
		return false, errors.New("database collection error")
	}
	_, err := coll.ReadDocument(nil, key, location)
	if err != nil {
		return false, err
	}
	return true, nil
}
func (l *LocationDB) UpdateLocation(patch map[string]interface{}, key string) (bool, error) {
	coll, ok := l.coll["locationCollection"]
	if !ok {
		return false, errors.New("database collection error")
	}
	_, err := coll.UpdateDocument(nil, key, patch)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (l *LocationDB) HealthCheck() (bool, error) {
	_, err := l.db.Info(nil)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (l *LocationDB) InsertLocation(location *structs.Location) (bool, error) {
	coll, ok := l.coll["locationCollection"]
	if !ok {
		return false, errors.New("database collection error")
	}
	_, err := coll.CreateDocument(nil, location)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (l *LocationDB) ReadLocations(filters map[string]string) ([]structs.Location, error) {
	bindVars := map[string]interface{}{}
	for i, value := range filters {
		switch i {
		case constants.Device:
			bindVars["deviceId"] = value
		case constants.Speed:
			number, err := strconv.ParseFloat(value, 64)
			if err != nil {
				return nil, err
			}
			bindVars["velocity"] = number
		case constants.Contract:
			number, err := strconv.Atoi(value)
			if err != nil {
				return nil, err
			}
			bindVars["contractId"] = number
		case constants.Fence:
			number, err := strconv.Atoi(value)
			if err != nil {
				return nil, err
			}
			bindVars[constants.Fence] = number
		case constants.Driver:
			bindVars[constants.Driver] = value
		case constants.InitialDate:
			layout := checkDateLayout(value)
			dateTime, err := time.Parse(layout, value)
			if err != nil {
				return nil, err
			}
			bindVars[constants.InitialDate] = dateTime
		case constants.FinalDate:
			layout := checkDateLayout(value)
			dateTime, err := time.Parse(layout, value)
			if err != nil {
				return nil, err
			}
			bindVars[constants.FinalDate] = dateTime
		}
	}
	query := "FOR l IN locationCollection FILTER "
	glueStr := ""
	for key := range bindVars {
		switch key {
		case constants.Fence:
			query += fmt.Sprintf("%s @%s IN l.%ss[*].fenceId", glueStr, key, key)
			break
		case constants.Driver:
			query += fmt.Sprintf("%sl.%s.id == @%s", glueStr, key, key)
			break
		case constants.InitialDate:
			query += fmt.Sprintf("%sl.locationDate >= @%s", glueStr, key)
			break
		case constants.FinalDate:
			query += fmt.Sprintf("%sl.locationDate <= @%s", glueStr, key)
			break
		default:
			query += fmt.Sprintf("%sl.%s == @%s", glueStr, key, key)
			break
		}
		glueStr = " AND "
	}
	query += " SORT l.locationDate DESC RETURN l"
	cursor, err := l.db.Query(nil, query, bindVars)
	if err != nil {
		return []structs.Location{}, errors.New("internal error")
	}
	var object []structs.Location
	for cursor.HasMore() == true {
		loc := structs.Location{}
		cursor.ReadDocument(nil, &loc)
		object = append(object, loc)
	}
	defer cursor.Close()
	return object, nil
}

func (l *LocationDB) ReadLocationWithFilters(location structs.Location) (structs.Location, error) {
	bindVars := map[string]interface{}{
		"_key":                  location.Key,
		"deviceId":              location.DeviceId,
		"systemType":            location.SystemType,
		"system":                location.System,
		"source":                location.Source,
		"creationDate":          location.CreationDate,
		"locationDate":          location.LocationDate,
		"coordinates":           location.Coordinates,
		"velocity":              location.Velocity,
		"precision":             location.Precision,
		"userId":                location.UserId,
		"userDescription":       location.UserDescription,
		"contractId":            location.ContractId,
		"deviceUniqueId":        location.DeviceUniqueId,
		"userTypeId":            location.UserTypeId,
		"userTypeDescription":   location.UserTypeDescription,
		"contractCode":          location.ContractCode,
		"enterpriseId":          location.EnterpriseId,
		"enterpriseDescription": location.EnterpriseDescription,
		"systemId":              location.SystemId,
		"systemDescription":     location.SystemDescription,
		"systemTypeId":          location.SystemTypeId,
		"systemTypeDescription": location.SystemTypeDescription,
		"timeLapse":             location.TimeLapse,
		"precisionQualifier":    location.PrecisionQualifier,
		"reason":                location.Reason,
		"version":               location.Version,
		"odometer":              location.Odometer,
		"horimetre":             location.HourMeter,
		"batteryLevel":          location.BatteryLevel,
		"extraFields":           location.ExtraFields,
		"integrator":            location.Integrator,
		"systemIntegrator":      location.SystemIntegrator,
	}
	query := "FOR l IN locationCollection FILTER "
	glueStr := ""
	for key := range bindVars {
		query += fmt.Sprintf("%sl.%s == @%s", glueStr, key, key)
		glueStr = " AND "
	}
	query += " RETURN l"
	cursor, err := l.db.Query(nil, query, bindVars)
	if err != nil {
		return structs.Location{}, errors.New("internal error")
	}
	object := structs.Location{}
	_, err = cursor.ReadDocument(nil, &object)
	if err != nil {
		return structs.Location{}, err
	}
	defer cursor.Close()
	return object, nil
}

func checkDateLayout(value string) string {
	var layout string
	if value[len(value)-1] == 'Z' {
		layout = "2006-01-02 15:04:05Z"
	} else {
		layout = "2006-01-02 15:04:05"
	}
	return layout
}
