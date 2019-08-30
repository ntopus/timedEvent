package config

type ConfigData struct {
	LogLevel       int           `json:"logLevel"`
	ServerPort     int           `json:"serverPort"`
	PoolTime       int           `json:"poolTime"`
	ControlTime    int           `json:"controlTime"`
	ExpirationTime int           `json:"expirationTime"`
	DataBase       ConfigDB      `json:"dataBase"`
	PublishQueue   []ConfigQueue `json:"loglevel"`
}

type ConfigQueue struct {
	ServerHost     string `json:"serverHost"`
	ServerVHost    string `json:"serverVHost"`
	ServerPort     string `json:"serverPort"`
	ServerUser     string `json:"serverUser"`
	ServerPassword string `json:"serverPassword"`
	QueueName      string `json:"queueName"`
}

type ConfigDB struct {
	ServerHost     string `json:"serverHost"`
	ServerPort     string `json:"serverPort"`
	ServerUser     string `json:"serverUser"`
	ServerPassword string `json:"serverPassword"`
	DbName         string `json:"dbName"`
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
