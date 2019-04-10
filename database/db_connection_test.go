package database

import (
	"fmt"
	"github.com/arangodb/go-driver"
	"github.com/arangodb/go-driver/http"
	"github.com/ivanmeca/timedEvent/config"
	"github.com/onsi/gomega"
	"testing"
)

func TestConnection(test *testing.T) {
	gomega.RegisterTestingT(test)

	fmt.Println("Trying to connect to a database server")

	connArgs := config.ConfigData{
		DataBase: config.ConfigDB{
			ServerHost:     "http://localhost",
			ServerPort:     "8529",
			ServerUser:     "root",
			ServerPassword: "rootpass1",
			DbName:         "testDb",
		}}
	_, err := NewClientDB(&connArgs.DataBase)
	gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
}

func TestLibConnection(test *testing.T) {
	gomega.RegisterTestingT(test)

	conn, err := http.NewConnection(http.ConnectionConfig{
		Endpoints: []string{"http://localhost:8529"},
	})
	gomega.Expect(err).ShouldNot(gomega.HaveOccurred())

	c, err := driver.NewClient(driver.ClientConfig{
		Connection:     conn,
		Authentication: driver.BasicAuthentication("root", "rootpass"),
	})
	gomega.Expect(err).ShouldNot(gomega.HaveOccurred())

	exist, err := c.DatabaseExists(nil, "dbTeste")
	gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
	gomega.Expect(exist).Should(gomega.BeFalse())

	_, err = c.CreateDatabase(nil, "dbTeste", nil)
	gomega.Expect(err).ShouldNot(gomega.HaveOccurred())

	//db, err := c.Database(nil, "_system")
	//if err != nil {
	//	// Handle error
	//}
	//
	//// Open "books" collection
	//col, err := db.Collection(nil, "books")
	//if err != nil {
	//	// Handle error
	//}

	//// Create document
	//book := Book{
	//	Title:   "ArangoDB Cookbook",
	//	NoPages: 257,
	//}
	//meta, err := col.CreateDocument(nil, book)
	//if err != nil {
	//	// Handle error
	//}
	//fmt.Printf("Created document in collection '%s' in database '%s'\n", col.Name(), db.Name())

}
