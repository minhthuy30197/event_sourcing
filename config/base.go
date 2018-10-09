package config

import (
	"log"
	"os"

	"git.hocngay.com/hocngay/gin-template/helper"
)

type Config struct {
	ServiceName string `json:"service_name"`

	Database struct {
		User     string `json:"user"`
		Password string `json:"password"`
		Database string `json:"database"`
		Address  string `json:"address"`
	} `json:"database"`

	EventDatabase struct {
		User     string `json:"user"`
		Password string `json:"password"`
		Database string `json:"database"`
		Address  string `json:"address"`
	} `json:"event_database"`
}

func SetupConfig() Config {
	var conf Config

	// Đọc file config.local.json
	configFile, err := os.Open("config.local.json")
	if err != nil {
		// Nếu không có file config.local.json thì đọc file config.default.json
		configFile, err = os.Open("config.default.json")
		if err != nil {
			panic(err)
		}
		defer configFile.Close()
	}
	defer configFile.Close()

	// Parse dữ liệu JSON và bind vào Controller
	err = helper.DecodeDataFromJsonFile(configFile, &conf)
	if err != nil {
		log.Println("Không đọc được file config.")
		panic(err)
	}

	return conf
}
