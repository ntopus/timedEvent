package database

import (
	"github.com/arangodb/go-driver"
	"github.com/arangodb/go-driver/http"
)

type DatabaseConfigurationReader interface {
	GetServerHost() string
	GetServerPort() string
	GetServerUser() string
	GetServerPassword() string
	GetDbName() string
}

func NewClientDB(configuration DatabaseConfigurationReader) (*driver.Client, error) {
	conn, err := http.NewConnection(http.ConnectionConfig{
		Endpoints: []string{configuration.GetServerHost() + ":" + configuration.GetServerPort()},
	})
	if err != nil {
		return nil, err
	}
	c, err := driver.NewClient(driver.ClientConfig{
		Connection:     conn,
		Authentication: driver.BasicAuthentication(configuration.GetServerUser(), configuration.GetServerPassword()),
	})
	if err != nil {
		return nil, err
	}
	return &c, nil
}
