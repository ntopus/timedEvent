package database

import (
	"github.com/arangodb/go-driver"
	"github.com/arangodb/go-driver/http"
)

type ArangoDBConnection struct {
	ServerHost     string
	ServerPort     string
	ServerUser     string
	ServerPassword string
	DbName         string
}

func NewClientDB(configuration *ArangoDBConnection) (*driver.Client, error) {
	conn, err := http.NewConnection(http.ConnectionConfig{
		Endpoints: []string{configuration.ServerHost + ":" + configuration.ServerPort},
	})
	if err != nil {
		return nil, err
	}
	c, err := driver.NewClient(driver.ClientConfig{
		Connection:     conn,
		Authentication: driver.BasicAuthentication(configuration.ServerUser, configuration.ServerPassword),
	})
	if err != nil {
		return nil, err
	}
	return &c, nil
}
