package arangoDB

import (
	"fmt"
	"github.com/ivanmeca/timedEvent/database"
	"github.com/ivanmeca/timedEvent/database/data_types"
	"github.com/onsi/gomega"
	"testing"
)

func TestReadDocuments(test *testing.T) {
	gomega.RegisterTestingT(test)
	fmt.Println("Trying to connect to a read collection")
	coll := getTestCollectionInstance("testeCollection")

	var item data_types.EventEntry
	var List []data_types.EventEntry
	err := coll.Read(nil, item, List)
	gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
	fmt.Println(List)
}

func getTestCollectionInstance(collName string) database.CollectionManagment {

	DBClient, err := NewDBClient(GetTestDatabase())
	gomega.Expect(err).ShouldNot(gomega.HaveOccurred())

	db, err := DBClient.GetDatabase("TestDB", true)
	gomega.Expect(err).ShouldNot(gomega.HaveOccurred())

	coll, err := db.GetCollection("TesteColl")

	return coll
}
