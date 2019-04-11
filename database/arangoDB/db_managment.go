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

func (dbm *ArangoDbManager) CreateCollection(collectionName string) (bool, error) {
	_, err := dbm.db.CreateCollection(nil, collectionName, nil)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (dbm *ArangoDbManager) DropCollection(collectionName string) (bool, error) {
	coll, err := dbm.db.Collection(nil, collectionName)
	if err != nil {
		return false, err
	}
	err = coll.Remove(nil)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (dbm *ArangoDbManager) GetCollection(collectionName string) (database.CollectionManagment, error) {
	coll, err := dbm.db.Collection(nil, collectionName)
	if err != nil {
		return nil, err
	}
	return &ArangoDbCollection{
		db:               dbm.db,
		collection:       collectionName,
		collectionDriver: coll,
	}, nil
}

func (dbm *ArangoDbManager) Drop() (bool, error) {
	err := dbm.db.Remove(nil)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (dbm *ArangoDbManager) HealthCheck() (bool, error) {
	_, err := dbm.db.Info(nil)
	if err != nil {
		return false, err
	}
	return true, nil
}
