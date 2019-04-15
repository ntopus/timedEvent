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
	fmt.Println("Trying to a read collection")
	coll := getTestCollectionInstance("testeCollection")

	list, err := coll.Read(map[string]interface{}{})
	gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
	fmt.Println(list)
}

func TestInsertDocument(test *testing.T) {
	gomega.RegisterTestingT(test)
	fmt.Println("Trying insert into read collection")
	coll := getTestCollectionInstance("testeCollection")

	for i := 0; i < 10; i++ {
		data := fmt.Sprintf(`"Teste data %d"`, i)
		event, err := data_types.NewCloudEventJsonV02("TestEvent", []byte(data), nil)
		gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
		ok, err := coll.Insert(event)
		gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
		gomega.Expect(ok).Should(gomega.BeTrue())
	}

}

func TestReadCollection(test *testing.T) {
	gomega.RegisterTestingT(test)
	fmt.Println("Trying to a read collection")
	coll := getTestCollectionInstance("testeCollection")

	list, err := coll.Read(map[string]interface{}{})
	gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
	fmt.Println(list)
}

func getTestCollectionInstance(collName string) database.CollectionManagment {

	DBClient, err := NewDBClient(GetTestDatabase())
	gomega.Expect(err).ShouldNot(gomega.HaveOccurred())

	db, err := DBClient.GetDatabase("TestDB", true)
	gomega.Expect(err).ShouldNot(gomega.HaveOccurred())

	coll, err := db.GetCollection("TesteColl")

	if err != nil {
		ok, err := db.CreateCollection("TesteColl")
		gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
		gomega.Expect(ok).Should(gomega.BeTrue())
	}

	return coll
}
