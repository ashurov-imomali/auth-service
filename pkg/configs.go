package pkg

import (
	"encoding/json"
	"os"
)

var Sms2FA bool

func GetConfigs() (*Config, error) {
	var conf Config
	bytes, err := os.ReadFile("./config/configs.json")
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(bytes, &conf); err != nil {
		return nil, err
	}
	Sms2FA = conf.Sms2FA
	return &conf, nil
}
