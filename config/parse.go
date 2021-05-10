package config

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/kataras/golog"
)

type Config struct {
	Env    string
	Server Server
}

type Server struct {
	Address string
}

var Props Config
var confPath string

func Parse() bool {
	confPath = os.Getenv("CONFIG")
	if confPath == "" {
		golog.Info("no config path provided, parsing conf from <PROJECT_DIR/>service.json")
		confPath = "service.json"
	}
	confData, err := ioutil.ReadFile(confPath)
	if err != nil {
		golog.Error(err)
		return false
	}
	err = json.Unmarshal(confData, &Props)
	if err != nil {
		golog.Error(err)
		return false
	}
	err = validate()
	if err != nil {
		golog.Error(err)
		return false
	}
	return true
}

func validate() error {
	return nil
}
