package arangoDB

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

	_, err := NewDBClient(GetTestDatabase())
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
	gomega.Expect(exist).Should(gomega.BeFalse())

	_, err = c.CreateDatabase(nil, "dbTeste", nil)
	gomega.Expect(err).ShouldNot(gomega.HaveOccurred())

	db, err := c.Database(ctx, "dbTeste")
	gomega.Expect(err).ShouldNot(gomega.HaveOccurred())

	err = db.Remove(ctx)
	gomega.Expect(err).ShouldNot(gomega.HaveOccurred())

	exist, err = c.DatabaseExists(ctx, "dbTeste")
	gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
	gomega.Expect(exist).Should(gomega.BeFalse())

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
