package config

import (
	"encoding/json"
	"github.com/micro/go-config/source/env"
	"github.com/onsi/gomega"
	"os"
	"testing"
)

func TestEnv_Read(t *testing.T) {

	expected := map[string]map[string]string{
		"mongodb": {
			"server": "fleet.db.interno.ntopus.com.br",
		},
		//"server": {
		//	"port": "8081",
		//},
	}

	os.Setenv("MONGODB_SERVER", "fleet.db.interno.ntopus.com.br")
	//os.Setenv("SERVER_PORT", "8081")

	source := env.NewSource()
	c, err := source.Read()
	if err != nil {
		t.Error(err)
	}

	var actual map[string]interface{}
	if err := json.Unmarshal(c.Data, &actual); err != nil {
		t.Error(err)
	}

	actualDB := actual["mongodb"].(map[string]interface{})

	for k, v := range expected["mongodb"] {
		a := actualDB[k]

		if a != v {
			t.Errorf("expected %v got %v", v, a)
		}
	}
}

func TestInitConfig(t *testing.T) {
	gomega.RegisterTestingT(t)
	os.Setenv("MONGODB_SERVER", "fleet.db.interno.ntopus.com.br")
	InitConfig()
	server := GetDatabaseHost()
	gomega.Expect(server).To(gomega.Equal("fleet.db.interno.ntopus.com.br"))
}
