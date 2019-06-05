package collection_managment

import (
	"github.com/ivanmeca/timedEvent/application/modules/config"
	"github.com/ivanmeca/timedEvent/application/modules/database"
	"github.com/ivanmeca/timedEvent/application/modules/database/arangoDB"
	"github.com/ivanmeca/timedEvent/application/modules/routes"
	"net/http"
	"regexp"
	"strings"
	"sync"
)

var once sync.Once
var DBInstance database.DataBaseManagment

func GetDBSession() database.DataBaseManagment {
	once.Do(func() {
		var err error

		dbConn, err := arangoDB.NewDBClient(GetTestDatabase())
		if err != nil {
			panic(err)
		}
		DBInstance, err = dbConn.GetDatabase("TestDB", false)
		if err != nil {
			panic(err)
		}
	})
	return DBInstance
}

func DefaultErrorHandler(err error) routes.JsendMessage {
	errMsg := routes.JsendMessage{}
	if strings.Contains(err.Error(), "unique constraint violated ") {
		dupKey := regexp.MustCompile("conflicting key:\\s*([\\w\\W]+)")
		matches := dupKey.FindStringSubmatch(err.Error())
		if len(matches) > 0 {
			errMsg.SetMessage("duplicated id " + matches[1])
		} else {
			errMsg.SetMessage(err.Error())
		}
		errMsg.SetStatus(http.StatusForbidden)
		return errMsg
	}
	errMsg.SetStatus(http.StatusInternalServerError)
	errMsg.SetMessage(err.Error())
	return errMsg
}

func GetTestDatabase() *config.ConfigDB {
	return &config.ConfigDB{
		ServerHost:     "http://localhost",
		ServerPort:     "8529",
		ServerUser:     "testUser",
		ServerPassword: "123456",
		DbName:         "testDb",
	}
}
