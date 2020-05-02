package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

//DBSound contains info about each sound
type DBSound struct {
	id       string   //obselete?
	Filepath string   `json:"filepath"`
	Command  string   `json:"command"`
	Enabled  string   `json:"enabled"`
	loaded   string   //obselete?
	buffer   [][]byte //Obselete?
}

func loadSoundsFromJSON() {
	//Check if file exists, if not create it and terminate the program.
	if _, err := os.Stat("./sounds.json"); os.IsNotExist(err) {
		log.Println("sounds.json did not exist and have been created.")
		file, _ := json.MarshalIndent(DBSoundList, "", " ")

		_ = ioutil.WriteFile("./sounds.json", file, 0644)
	}

	file, err := os.Open("./sounds.json")
	if err != nil {
		log.Println(err)
	}
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&DBSoundList)
	if err != nil {
		log.Println(err)
	}
}

func saveSoundsToJSON() {
	file, _ := json.MarshalIndent(DBSoundList, "", " ")

	_ = ioutil.WriteFile("./sounds.json", file, 0644)

}
