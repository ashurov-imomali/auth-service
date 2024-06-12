package pkg

import (
	"gopkg.in/yaml.v3"
	"os"
)

var Params = &TFAParams{}

func GetConfigs() (*Config, error) {
	var conf Config
	bytes, err := os.ReadFile("./config/configs.yaml")
	if err != nil {
		return nil, err
	}
	if err := yaml.Unmarshal(bytes, &conf); err != nil {
		return nil, err
	}
	Params = conf.TFAParams
	return &conf, nil
}
