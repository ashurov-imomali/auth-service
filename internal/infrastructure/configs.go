package infrastructure

import (
	"encoding/json"
	"main/pkg"
	"os"
)

func GetConfigs() (*pkg.Config, error) {
	var conf pkg.Config
	bytes, err := os.ReadFile("./config/configs.json")
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(bytes, &conf); err != nil {
		return nil, err
	}
	return &conf, nil
}
