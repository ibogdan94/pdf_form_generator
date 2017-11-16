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

var Config = Properties{}

func ParseJSONConfig() (props Properties, err error) {
	if Config.Port != "" {
		return Config, nil
	}

	pwd, err := os.Getwd()

	if err != nil {
		return Config, err
	}

	payload, err := ioutil.ReadFile(pwd + "/config.json")

	if err != nil {
		return Config, err
	}

	if err := json.Unmarshal(payload, &Config); err != nil {
		return Config, err
	}

	return Config, nil
}
