package arangoDB

import (
	"errors"
	"github.com/arangodb/go-driver"
	"github.com/arangodb/go-driver/http"
	"github.com/ivanmeca/timedEvent/database"
)

func NewDBClient(configuration database.DatabaseConfigurationReader) (database.DataBaseConnector, error) {
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

	return &ArangoDBConnector{
		conn:   conn,
		client: c,
	}, nil
}

type ArangoDBConnector struct {
	conn   driver.Connection
	client driver.Client
}

func (db *ArangoDBConnector) GetDatabase(databaseName string, createIfNotExists bool) (database.DataBaseManagment, error) {
	exist, err := db.client.DatabaseExists(nil, databaseName)
	if err != nil {
		return nil, errors.New("could not create database: " + err.Error())
	}
	if !exist && createIfNotExists {
		_, err := db.client.CreateDatabase(nil, databaseName, nil)
		if err != nil {
			return nil, errors.New("could not create database: " + err.Error())
		}
	}
	d, err := db.client.Database(nil, databaseName)
	if err != nil {
		return nil, errors.New("could not create database: " + err.Error())
	}
	return &ArangoDbManager{
		client: db.client,
		conn:   db.conn,
		db:     d,
	}, nil
}
