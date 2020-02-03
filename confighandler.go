package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

type Config struct {
	ConnectionString string
	DiscToken        string
	CommandText      string
}

var config Config

func ConfigInit() {
	//Check if file exists anton_config.json
	if _, err := os.Stat("./tpudiscordbot.json"); os.IsNotExist(err) {
		log.Println("tpudiscordbot.json did not exist and have been created. Please fill in the fields")
		file, _ := json.MarshalIndent(config, "", " ")

		_ = ioutil.WriteFile("./tpudiscordbot.json", file, 0644)
		os.Exit(1)
	}

	//Load config.
	file, err := os.Open("./tpudiscordbot.json")
	if err != nil {
		log.Println(err)
	}
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		log.Println(err)
	}

	//log.Println(config)
}
