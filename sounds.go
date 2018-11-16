package main

import (
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/bwmarrin/discordgo"
)

//DBSound contains info about each sound
type DBSound struct {
	id       string
	filepath string //Backwards compatibility
	command  string
	enabled  string
	loaded   string
	buffer   [][]byte
	noplays  string
}

//DBSoundList is an array of DBSound items
var DBSoundList []DBSound

func playSound(s *discordgo.Session, guildID, channelID string, command string) (err error) {

	var index = -1
	for i := range DBSoundList {
		if DBSoundList[i].command == command {
			index = i
		}
	}

	if index == -1 {
		return
	}

	if DBSoundList[index].enabled == "0" {
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
	vc.Speaking(true)

	// Send the buffer data.
	for _, buff := range DBSoundList[index].buffer {
		if currentlyPlaying[guildID] == "" {
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
	loadFromDB()

	for i := range DBSoundList {
		if DBSoundList[i].enabled != "0" && DBSoundList[i].loaded != "1" {
			loadSound(DBSoundList[i].filepath, DBSoundList[i].command)
			DBSoundList[i].loaded = "1"
			log.Println("Loaded " + DBSoundList[i].command)
		}
	}
}

func addToList(obj DBSound) {

	var index int
	index = -1
	for i := range DBSoundList {
		if obj.id == DBSoundList[i].id {
			index = i
		}
	}

	if index != -1 {
		DBSoundList[index].enabled = obj.enabled
		DBSoundList[index].command = obj.command

		var index2 int
		index2 = -1
		for i := range DBSoundList {
			if obj.command == DBSoundList[i].command {
				index2 = i
			}
		}

		if obj.enabled == "1" {
			DBSoundList[index2].enabled = "1"
		} else {
			DBSoundList[index2].enabled = "0"
		}

		if index2 == -1 {
			log.Println("Holy fucking moly. Something went wrong")
		}

	} else if obj.enabled == "1" {
		//If it dosent exist add it.
		DBSoundList = append(DBSoundList, obj)
		log.Println("added " + obj.command + " to the list")
	}

}

// loadSound attempts to load an encoded sound file from disk.
func loadSound(path string, command string) error {

	var index int
	index = -1
	for i := range DBSoundList {
		if DBSoundList[i].filepath == path {
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
