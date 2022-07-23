package main

import (
	jsoniter "github.com/hezhis/go"
	"io/ioutil"
)

var config *Config

type Config struct {
	ClientPath string `json:"clientPath"`
	ClientType string `json:"clientType"`
	ServerPath string `json:"serverPath"`
	ServerType string `json:"serverType"`
}

func loadConfig() error {
	data, err := ioutil.ReadFile("./config.json")
	if nil != err {
		return err
	}

	return jsoniter.Unmarshal(data, &config)
}
