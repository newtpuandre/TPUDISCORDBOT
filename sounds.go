package main

import (
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

//DBSoundList is an array of DBSound items
var DBSoundList []DBSound

func playSound(s *discordgo.Session, guildID, channelID string, command string) (err error) {

	var index = -1
	for i := range DBSoundList {
		if DBSoundList[i].Command == command {
			index = i
		}
	}

	if index == -1 {
		return
	}

	if DBSoundList[index].Enabled == "0" {
		s.ChannelMessageSend(channelID, "That command is disabled.")
		return
	}

	// Join the provided voice channel.
	vc, err := s.ChannelVoiceJoin(guildID, channelID, false, true)
	if err != nil {
		return err
	}

	// Sleep for a specified amount of time before playing the sound
	time.Sleep(250 * time.Millisecond)

	// Start speaking.
	err = vc.Speaking(true)
	if err != nil {
		vc.Disconnect()
	}

	// Send the buffer data.
	for _, buff := range DBSoundList[index].buffer {
		if currentlyPlaying[guildID] == "" { //If !stop command is sent
			break
		}
		vc.OpusSend <- buff
	}

	// Stop speaking
	vc.Speaking(false)

	// Sleep for a specificed amount of time before ending.
	time.Sleep(250 * time.Millisecond)

	// Disconnect from the provided voice channel.
	vc.Disconnect()

	return nil
}

func loadFromList() {

	var files []string

	root := "./sounds/"
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		files = append(files, path)
		return nil
	})
	if err != nil {
		panic(err)
	}

	for id, file := range files {
		if strings.Contains(file, "/sounds/") {
			continue //Skip loop iteration if the file is the root directory.
		}
		var tempDBSound DBSound
		tempDBSound.id = string(id)
		tempDBSound.Enabled = "1"
		tempDBSound.Filepath = file

		tempCommandString := strings.TrimLeft(file, "sounds/") //Removes the path /sounds/
		if strings.Contains(tempCommandString, "\\") {
			tempCommandString = strings.TrimLeft(tempCommandString, "\\")
		}
		tempCommandString = strings.TrimSuffix(tempCommandString, ".dca") //Removes the file extension

		if strings.Contains(tempCommandString, "!") {
			//tempCommandString = tempCommandString[0:] //Trim of the !
			tempCommandString = strings.TrimLeft(tempCommandString, "!")
		}

		tempDBSound.Command = tempCommandString

		addToList(tempDBSound)
	}

	for i := range DBSoundList {
		if DBSoundList[i].loaded != "1" && DBSoundList[i].Enabled != "0" {
			loadSound(DBSoundList[i].Filepath, DBSoundList[i].Command)
			DBSoundList[i].loaded = "1"
			log.Println("Loaded " + DBSoundList[i].Command)
		}
	}
}

func addToList(obj DBSound) {

	for i := range DBSoundList {
		if obj.Command == DBSoundList[i].Command {
			return
		}
	}

	DBSoundList = append(DBSoundList, obj)
	log.Println("added " + obj.Command + " to the list")

}

// loadSound attempts to load an encoded sound file from disk.
func loadSound(path string, command string) error {

	var index int
	index = -1
	for i := range DBSoundList {
		if DBSoundList[i].Filepath == path {
			index = i
		}
	}

	var tempBuffer = make([][]byte, 0)

	file, err := os.Open(path)
	if err != nil {
		fmt.Println("Error opening dca file :", err)
		return err
	}

	var opuslen int16

	for {
		// Read opus frame length from dca file.
		err = binary.Read(file, binary.LittleEndian, &opuslen)

		// If this is the end of the file, just return.
		if err == io.EOF || err == io.ErrUnexpectedEOF {
			err := file.Close()
			if err != nil {
				return err
			}
			DBSoundList[index].buffer = tempBuffer
			return nil
		}

		if err != nil {
			fmt.Println("Error reading from dca file :", err)
			return err
		}

		// Read encoded pcm from dca file.
		InBuf := make([]byte, opuslen)
		err = binary.Read(file, binary.LittleEndian, &InBuf)

		// Should not be any end of file errors
		if err != nil {
			fmt.Println("Error reading from dca file :", err)
			return err
		}

		// Append encoded pcm data to the buffer.

		tempBuffer = append(tempBuffer, InBuf)

	}
}

func soundExist(command string) bool {
	var index int
	index = -1
	for i := range DBSoundList {
		if DBSoundList[i].Command == command {
			index = i
		}
	}

	if index != -1 {
		return true
	}

	return false

}
