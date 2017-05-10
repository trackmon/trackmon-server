package main

import (
	"fmt"
	"io/ioutil"
)

type Configuration struct {
	ListeningAddress string
	DatabaseAddress  string
	DatabasePassword string
	AutoUpdateChecker bool
}

func RecreateConfig(path string) {
	return
}

func CreateConfig() {
	var Config Configuration

	// Standard config
	Config.ListeningAddress = ":80"
	Config.DatabaseAddress = "localhost"
	Config.DatabasePassword = ""
	Config.AutoUpdateChecker = true

	ByteJsonConfig, err := toprettyjson(Config)
	if err != nil {
		panic(err)
	}
	fmt.Println("Writing configuration to ./trackmonserv.conf")
	err = ioutil.WriteFile("./trackmonserv.conf", ByteJsonConfig, 0644)
	if err != nil {
		panic(err)
	}
}
