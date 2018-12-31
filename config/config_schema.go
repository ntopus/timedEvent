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
