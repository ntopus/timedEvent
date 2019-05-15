package arangoDB

import (
	"github.com/arangodb/go-driver"
	"github.com/ivanmeca/timedEvent/database"
)

type Manager struct {
	conn   driver.Connection
	client driver.Client
	db     driver.Database
}

func (dbm *Manager) CreateCollection(collectionName string) (bool, error) {
	_, err := dbm.db.CreateCollection(nil, collectionName, nil)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (dbm *Manager) DropCollection(collectionName string) (bool, error) {
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

func (dbm *Manager) GetCollection(collectionName string) (database.CollectionManagment, error) {
	coll, err := dbm.db.Collection(nil, collectionName)
	if err != nil {
		return nil, err
	}
	return &Collection{
		db:               dbm.db,
		collection:       collectionName,
		collectionDriver: coll,
	}, nil
}

func (dbm *Manager) Drop() (bool, error) {
	err := dbm.db.Remove(nil)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (dbm *Manager) HealthCheck() (bool, error) {
	_, err := dbm.db.Info(nil)
	if err != nil {
		return false, err
	}
	return true, nil
}
