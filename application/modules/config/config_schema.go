package config

type ConfigData struct {
	LogLevel       int
	PoolTime       int
	ControlTime    int
	ExpirationTime int
	DataBase       ConfigDB
	PublishQueue   []ConfigQueue
}

type ConfigQueue struct {
	ServerHost     string
	ServerVHost    string
	ServerPort     string
	ServerUser     string
	ServerPassword string
	QueueName      string
}

type ConfigDB struct {
	ServerHost     string
	ServerPort     string
	ServerUser     string
	ServerPassword string
	DbName         string
}
