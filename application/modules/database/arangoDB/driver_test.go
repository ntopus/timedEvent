package arangoDB

import (
	"errors"
	"fmt"
	"github.com/arangodb/go-driver"
	"github.com/arangodb/go-driver/http"
	"github.com/ivanmeca/timedEvent/application/modules/database/data_types"
	"github.com/onsi/gomega"
	"sync"
	"testing"
	"time"
)

func getDBConn() (driver.Client, error) {
	conn, err := http.NewConnection(http.ConnectionConfig{
		Endpoints: []string{"tcp://127.0.0.1:8529"},
	})
	if err != nil {
		return nil, err
	}
	c, err := driver.NewClient(driver.ClientConfig{
		Connection:     conn,
		Authentication: driver.BasicAuthentication("testUser", "123456"),
	})
	if err != nil {
		return nil, err
	}
	return c, nil
}

func getDB(conn driver.Client, databaseName string) (driver.Database, error) {
	exist, err := conn.DatabaseExists(nil, databaseName)
	if err != nil {
		return nil, errors.New("could not create database: " + err.Error())
	}
	if !exist {
		_, err := conn.CreateDatabase(nil, databaseName, nil)
		if err != nil {
			return nil, errors.New("could not create database: " + err.Error())
		}
	}
	d, err := conn.Database(nil, databaseName)
	if err != nil {
		return nil, errors.New("could not create database: " + err.Error())
	}
	return d, nil
}

func getCollection(Db driver.Database, collName string) error {

}

func TestDriverAsyncInsertMultipleDocuments(test *testing.T) {
	gomega.RegisterTestingT(test)
	fmt.Println("Trying insert into read collection")
	coll := getTestCollectionInstance("testeCollection")
	horaAtual := time.Now().UTC()
	wg := sync.WaitGroup{}
	for i := 0; i < 10000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			data := fmt.Sprintf(`"Teste data %d"`, i)
			event, err := data_types.NewArangoCloudEventV02("TestEvent", data, nil)
			gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
			publishdate := horaAtual.Add(time.Duration(i*60) * time.Second).Format("2006-01-02 15:04:05Z")
			event.PublishDate = publishdate
			eventTime := horaAtual.AddDate(0, 0, i)
			err = event.SetTime(eventTime)
			gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
			timeInitRequest := time.Now()
			newDoc, err := coll.Insert(event)
			timeDiff := time.Now().Sub(timeInitRequest)
			gomega.Expect(timeDiff).To(gomega.BeNumerically("<", 9500*time.Millisecond))
			gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
			readDocument(newDoc.GetID())
		}()
	}
	wg.Wait()
}

func TestDriverInsertMultipleDocuments(test *testing.T) {
	gomega.RegisterTestingT(test)
	fmt.Println("Trying insert into read collection")
	coll := getTestCollectionInstance("testeCollection")
	horaAtual := time.Now().UTC()
	wg := sync.WaitGroup{}
	for i := 0; i < 10000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			data := fmt.Sprintf(`"Teste data %d"`, i)
			event, err := data_types.NewArangoCloudEventV02("TestEvent", data, nil)
			gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
			publishdate := horaAtual.Add(time.Duration(i*60) * time.Second).Format("2006-01-02 15:04:05Z")
			event.PublishDate = publishdate
			eventTime := horaAtual.AddDate(0, 0, i)
			err = event.SetTime(eventTime)
			gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
			timeInitRequest := time.Now()
			newDoc, err := coll.Insert(event)
			timeDiff := time.Now().Sub(timeInitRequest)
			gomega.Expect(timeDiff).To(gomega.BeNumerically("<", 9500*time.Millisecond))
			gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
			readDocument(newDoc.GetID())
		}()
	}
	wg.Wait()
}
