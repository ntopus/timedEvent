package database

import (
	"errors"
	"github.com/arangodb/go-driver"
	"github.com/arangodb/go-driver/http"
)

func NewClientDB(configuration *config.AppConfig) (*driver.Client, error) {
	baseUrl, ok := configuration.GetConfigParam(config.BaseUrl, "").(string)
	if !ok {
		return nil, errors.New("error getting BaseUrl from config.json")
	}
	port, ok := configuration.GetConfigParam(config.Port, "").(string)
	if !ok {
		return nil, errors.New("error getting Port from config.json")
	}
	conn, err := http.NewConnection(http.ConnectionConfig{
		Endpoints: []string{baseUrl + ":" + port},
	})
	if err != nil {
		return nil, err
	}
	user, ok := configuration.GetConfigParam(config.User, "").(string)
	if !ok {
		return nil, errors.New("error getting User from config.json")
	}
	password, ok := configuration.GetConfigParam(config.Password, "").(string)
	if !ok {
		return nil, errors.New("error getting Password from config.json")
	}
	c, err := driver.NewClient(driver.ClientConfig{
		Connection:     conn,
		Authentication: driver.BasicAuthentication(user, password),
	})
	if err != nil {
		return nil, err
	}
	return &c, nil
}
