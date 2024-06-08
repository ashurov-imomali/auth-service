package pkg

import (
	"encoding/json"
	"os"
)

var Params = &TFAParams{}

func GetConfigs() (*Config, error) {
	var conf Config
	bytes, err := os.ReadFile("./config/configs.json")
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(bytes, &conf); err != nil {
		return nil, err
	}
	Params = conf.TFAParams
	return &conf, nil
}
