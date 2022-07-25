package config

import (
	jsoniter "github.com/hezhis/go"
	"io/ioutil"
)

var config *Config

type Config struct {
	ClientPath    string `json:"clientPath"`
	ClientType    string `json:"clientType"`
	ServerPath    string `json:"serverPath"`
	ServerType    string `json:"serverType"`
	ExportDefault bool   `json:"exportDefault"` // 是否导出默认值
}

func GetPath(b bool) string {
	if b {
		return config.ServerPath
	}
	return config.ClientPath
}

func GetType(b bool) string {
	if b {
		return config.ServerType
	}
	return config.ClientType
}

func IsExportDefault() bool {
	return config.ExportDefault
}

func LoadConfig() error {
	data, err := ioutil.ReadFile("./config.json")
	if nil != err {
		return err
	}

	return jsoniter.Unmarshal(data, &config)
}
