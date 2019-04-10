package database

import (
	"fmt"
	"github.com/ivanmeca/timedEvent/config"
	"github.com/onsi/gomega"
	"testing"
)

func TestReadDocuments(test *testing.T) {
	gomega.RegisterTestingT(test)
	fmt.Println("Trying to connect to a read collection")
	db := getDbInstance()
	items, err := db.ReadCollection(eventsCollection, nil)
	gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
	fmt.Println(items)
}

func getDbInstance() Management {
	connArgs := config.ConfigData{
		DataBase: config.ConfigDB{
			ServerHost:     "http://127.0.0.1",
			ServerPort:     "8529",
			ServerUser:     "root",
			ServerPassword: "rootpass",
			DbName:         "testDb",
		}}
	db, err := NewDbManagement(&connArgs.DataBase)
	gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
	return db
}
