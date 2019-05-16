package arangoDB

import (
	"fmt"
	"github.com/ivanmeca/timedEvent/application/modules/database"
	"github.com/ivanmeca/timedEvent/application/modules/database/data_types"
	"github.com/onsi/gomega"
	"testing"
	"time"
)

func TestReadDocumentsWithFilter(test *testing.T) {
	gomega.RegisterTestingT(test)
	fmt.Println("Trying to a read collection")
	coll := getTestCollectionInstance("testeCollection")
	horaAtual := time.Now().AddDate(0, 0, 3)
	list, err := coll.Read([]database.AQLComparator{{Field: "Context.time", Comparator: ">=", Value: horaAtual}})
	gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
	fmt.Println(list)
}

func TestReadDocuments(test *testing.T) {
	gomega.RegisterTestingT(test)
	fmt.Println("Trying to a read collection")
	coll := getTestCollectionInstance("testeCollection")

	list, err := coll.Read(nil)
	gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
	fmt.Println(list)
}

func TestInsertDocument(test *testing.T) {
	gomega.RegisterTestingT(test)
	fmt.Println("Trying insert into read collection")
	coll := getTestCollectionInstance("testeCollection")

	horaAtual := time.Now()

	for i := 0; i < 10; i++ {
		data := fmt.Sprintf(`"Teste data %d"`, i)
		event, err := data_types.NewCloudEventV02("TestEvent", data, nil)
		gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
		eventTime := horaAtual.AddDate(0, 0, i)
		event.Context.SetTime(eventTime)
		ok, err := coll.Insert(event)
		gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
		gomega.Expect(ok).Should(gomega.BeTrue())
	}

}

func TestReadCollection(test *testing.T) {
	gomega.RegisterTestingT(test)
	fmt.Println("Trying to a read collection")
	coll := getTestCollectionInstance("testeCollection")

	list, err := coll.Read(nil)
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
