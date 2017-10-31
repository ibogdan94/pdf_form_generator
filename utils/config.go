package utils

import (
	"os"
	"io/ioutil"
	"encoding/json"
)

type Properties struct {
	Port string
	StaticPath string
	NodeModulesPath string
	TempPath string
	ServerPem string
	ServerKey string
}

var config = Properties{}

func ParseJSONConfig() (props Properties, err error) {
	if config.Port != "" {
		return config, nil
	}

	pwd, err := os.Getwd()

	if err != nil {
		return config, err
	}

	payload, err := ioutil.ReadFile(pwd + "/config.json")

	if err != nil {
		return config, err
	}

	if err := json.Unmarshal(payload, &config); err != nil {
		return config, err
	}


	return config, nil
}