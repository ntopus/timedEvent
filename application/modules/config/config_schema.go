package config

type ConfigData struct {
	LogLevel int
	DataBase ConfigDB
}

type ConfigDB struct {
	ServerHost     string
	ServerPort     string
	ServerUser     string
	ServerPassword string
	DbName         string
}

func (db *ConfigDB) GetServerHost() string {
	return db.ServerHost
}
func (db *ConfigDB) GetServerPort() string {
	return db.ServerPort
}
func (db *ConfigDB) GetServerUser() string {
	return db.ServerUser
}
func (db *ConfigDB) GetServerPassword() string {
	return db.ServerPassword
}
func (db *ConfigDB) GetDbName() string {
	return db.DbName
}
