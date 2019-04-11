package arangoDB

import (
	"github.com/arangodb/go-driver"
	"github.com/ivanmeca/timedEvent/database"
)

type ArangoDbManager struct {
	conn   driver.Connection
	client driver.Client
	db     driver.Database
}

func (dbm *ArangoDbManager) GetCollection(collectionName string) (database.CollectionManagment, error) {

	return &ArangoDbCollection{}, nil
}

func (dbm *ArangoDbManager) Drop() (bool, error) {
	panic("implement me")
}

func (dbm *ArangoDbManager) HealthCheck() (bool, error) {
	_, err := dbm.db.Info(nil)
	if err != nil {
		return false, err
	}
	return true, nil
}
