package config

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

type Configure struct {
	PackageName  string
	MainActivity string
	SDKPath      string
}

var configPath = "configure.json"
var configuration *Configure = &Configure{}

func InitConfig() {
	content, err := ioutil.ReadFile(configPath)
	if err != nil {
		log.Fatalln(err)
	}

	err = json.Unmarshal(content, configuration)
	if err != nil {
		log.Fatalln(err)
	}

	if len(configuration.MainActivity) <= 0 || len(configuration.PackageName) <= 0 || len(configuration.SDKPath) <= 0 {
		log.Fatalln("Configuration error:", configuration.MainActivity, configuration.PackageName, configuration.SDKPath)
	}
}

func GetPackageName() string {
	return configuration.PackageName
}

func GetMainActivity() string {
	return configuration.MainActivity
}

func GetSDKPath() string {
	return configuration.SDKPath
}
