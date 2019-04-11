package database

type DatabaseConfigurationReader interface {
	GetServerHost() string
	GetServerPort() string
	GetServerUser() string
	GetServerPassword() string
}

type DataBaseConnector interface {
	GetDatabase(databaseName string, createIfNotExists bool) (DataBaseManagment, error)
}

type DataBaseManagment interface {
	CreateCollection(collectionName string) (bool, error)
	DropCollection(collectionName string) (bool, error)
	GetCollection(collectionName string) (CollectionManagment, error)
	Drop() (bool, error)
	HealthCheck() (bool, error)
}

type CollectionManagment interface {
	Insert(item interface{}) (bool, error)
	DeleteItem(keyList []string) (bool, error)
	Update(patch map[string]interface{}, key string) (bool, error)
	Read(filters map[string]interface{}, items []interface{}) error
	ReadItem(key string, item interface{}) error
}
