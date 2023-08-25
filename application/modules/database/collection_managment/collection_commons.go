package collection_managment

import (
	"fmt"
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
		appConfig := config.GetConfig()
		dbConn, err := arangoDB.NewDBClient(&appConfig.DataBase)
		if err != nil {
			fmt.Println(err)
			panic(err)
		}
		DBInstance, err = dbConn.GetDatabase(appConfig.DataBase.DbName, true)
		if err != nil {
			fmt.Println(err)
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
		errMsg.SetStatus(http.StatusConflict)
		return errMsg
	}
	errMsg.SetStatus(http.StatusInternalServerError)
	errMsg.SetMessage(err.Error())
	return errMsg
}
