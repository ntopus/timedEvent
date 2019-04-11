package database

import (
	"context"
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
		Authentication: driver.BasicAuthentication("testUser", "123456"),
	})
	gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
	ctx := context.Background()

	exist, err := c.DatabaseExists(ctx, "dbTeste")
	gomega.Expect(err).ShouldNot(gomega.HaveOccurred())

	if !exist {
		_, err = c.CreateDatabase(nil, "dbTeste", nil)
		gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
	}

	db, err := c.Database(ctx, "dbTeste")
	gomega.Expect(err).ShouldNot(gomega.HaveOccurred())

	err = db.Remove(ctx)
	gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
}
