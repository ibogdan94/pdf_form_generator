package utils

import (
	"os"
	"io/ioutil"
	"encoding/json"
)

type Properties struct {
	Port            string `json:"port"`
	StaticPath      string `json:"static_path"`
	NodeModulesPath string `json:"node_modules_path"`
	TempPath        string `json:"temp_path"`
	ServerPem       string `json:"server_pem"`
	ServerKey       string `json:"server_key"`
	Env             string `json:"env"`
	LogFileName     string `json:"log_file_name"`
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
