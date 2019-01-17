package database

import (
	"github.com/globalsign/mgo"
	"github.com/ivanmeca/timedQueueService/config"
	"net/http"
	"regexp"
	"strings"
	"sync"
)

var dbSession *mgo.Session
var once sync.Once

const FleetDBName = "fleet-management"

func EnsureIndex() error {
	err := DriverEnsureIndex()
	if err != nil {
		return err
	}
	err = GsmEnsureIndex()
	if err != nil {
		return err
	}
	err = TrackingDeviceEnsureIndex()
	if err != nil {
		return err
	}
	err = VehicleEnsureIndex()
	if err != nil {
		return err
	}
	return nil
}

func GetDBSession() *mgo.Session {
	once.Do(func() {
		var err error
		server := config.GetConfig().DataBase.ServerHost
		dbSession, err = mgo.Dial(server)
		if err != nil {
			panic(err)
		}
		//c:= mgo.Credential{}
		//c.Username = "root"
		//c.Password = "example"
		//lError:=dbSession.Login(&c)
		//if lError!= nil {
		//	log.Fatal(err)
		//}
		//dbSession.SetMode(mgo.Monotonic, true)
	})
	dbSession.Refresh()
	return dbSession
}

func InsertIntoCollection(DB string, Collection string, data interface{}) (bool, error) {
	err := getCollection(DB, Collection).Insert(data)
	if err != nil {
		return false, err
	}
	return true, nil
}

func getCollection(DB string, Collection string) *mgo.Collection {
	return GetDBSession().DB(DB).C(Collection)
}

func isValidObjectId(s string) bool {
	m, err := regexp.MatchString("^[\\dabcdefABCDEF]{24}$", s)
	if err != nil {
		return false
	}
	return m
}

func parseQuery(query interface{}) interface{} {
	return query
	queryObj := data_types.NewQueryType(query)
	return queryObj.GetMongoQuery()
}
