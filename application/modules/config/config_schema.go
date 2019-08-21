package config

type ConfigData struct {
	LogLevel       int
	ServerPort     int
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

func (c *ConfigDB) GetServerHost() string {
	return c.ServerHost
}

func (c *ConfigDB) GetServerPort() string {
	return c.ServerPort
}

func (c *ConfigDB) GetServerUser() string {
	return c.ServerUser
}

func (c *ConfigDB) GetServerPassword() string {
	return c.ServerPassword
}
