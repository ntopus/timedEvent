package file_config

type AppConfig struct {
	Port        int    `json:"port"`
	TokenFolder string `json:"tokenFolder"`
}

var currentConfig *AppConfig

func init() {
	currentConfig = NewCurrentConfig()
	return
}

func NewCurrentConfig() *AppConfig {
	return &AppConfig{}
}

func GetConfig() *AppConfig {
	return currentConfig
}

func SetConfig(config *AppConfig) {
	currentConfig = config
}
