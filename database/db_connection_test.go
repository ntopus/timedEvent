package database

import (
	"fmt"
	"github.com/onsi/gomega"
	"testing"
)

func TestConnection(test *testing.T) {
	gomega.RegisterTestingT(test)

	fmt.Println("Trying to connect to a database server")

	connArgs := ArangoDBConnection{
		ServerHost:     "timedEvent.db.ivanmeca.com.br",
		ServerPort:     "8529",
		ServerUser:     "root",
		ServerPassword: " rootpass",
		DbName:         "testDb",
	}
	_, err := NewClientDB(&connArgs)
	gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
}
