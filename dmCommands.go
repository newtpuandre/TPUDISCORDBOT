package main

import (
	"log"
	"strings"

	"github.com/bwmarrin/discordgo"
)

//func uploadAudio
//This ^ function will take a file by DM and add it to the commands

//These command are only allowed for people who
//are given access

func enableCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	channel, err := s.UserChannelCreate(m.Author.ID)
	if err != nil {
		log.Println("Something went wrong whilst trying to create a DM, " + err.Error())
	}

	words := strings.Split(m.Content, " ")
	if len(words) > 1 {

		if !soundExist(words[1]) {
			_, err = s.ChannelMessageSend(channel.ID, "This command does not exist")
			if err != nil {
				log.Println("Could not send message to discord in manageAPIKey(), " + err.Error())
			}
			return
		}

		commandInput := words[1]
		err := DBenableCommand(commandInput)
		if err != nil {
			log.Println(err)
		}
	}

	log.Println("Enabled command:", words[1])

	_, err = s.ChannelMessageSend(channel.ID, "Enabled Command: "+words[1])
	if err != nil {
		log.Println("Could not send message to discord in manageAPIKey(), " + err.Error())
	}
	loadFromList()
}

func disableCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	channel, err := s.UserChannelCreate(m.Author.ID)
	if err != nil {
		log.Println("Something went wrong whilst trying to create a DM, " + err.Error())
	}

	words := strings.Split(m.Content, " ")
	if len(words) > 1 {

		if !soundExist(words[1]) {
			_, err = s.ChannelMessageSend(channel.ID, "This command does not exist")
			if err != nil {
				log.Println("Could not send message to discord in manageAPIKey(), " + err.Error())
			}
			return
		}

		commandInput := words[1]
		err := DBdisableCommand(commandInput)
		if err != nil {
			log.Println(err)
		}
	}

	log.Println("Disabled command:", words[1])

	_, err = s.ChannelMessageSend(channel.ID, "Disabled Command: "+words[1])
	if err != nil {
		log.Println("Could not send message to discord in manageAPIKey(), " + err.Error())
	}
	loadFromList()
}
