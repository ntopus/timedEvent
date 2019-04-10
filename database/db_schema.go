package database

import (
	"github.com/arangodb/go-driver"
)

func schema(client driver.Client, dbName string, collections []string) (driver.Database, map[string]driver.Collection, error) {
	db, err := checkDatabase(client, dbName)
	if err != nil {
		return nil, nil, err
	}
	collMap, err := checkCollections(db, collections)
	if err != nil {
		return nil, nil, err
	}
	return db, collMap, nil
}

func checkDatabase(client driver.Client, dbName string) (driver.Database, error) {
	exist, err := client.DatabaseExists(nil, dbName)
	if err != nil {
		return nil, err
	}
	var db driver.Database
	if !exist {
		db, err := client.CreateDatabase(nil, dbName, nil)
		if err != nil {
			return nil, err
		}
		return db, nil
	}
	db, err = client.Database(nil, dbName)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func checkCollections(db driver.Database, collections []string) (map[string]driver.Collection, error) {
	collMap := make(map[string]driver.Collection)
	for _, collection := range collections {
		exist, err := db.CollectionExists(nil, collection)
		if err != nil {
			return nil, err
		}
		if !exist {
			coll, err := db.CreateCollection(nil, collection, nil)
			if err != nil {
				return nil, err
			}
			collMap[collection] = coll
			continue
		}
		coll, err := db.Collection(nil, collection)
		if err != nil {
			return nil, err
		}
		collMap[collection] = coll
	}
	return collMap, nil
}
