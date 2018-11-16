package main

import (
	"log"
	"strings"

	"github.com/bwmarrin/discordgo"
)

var inUseServers map[string]string

var currentlyPlaying map[string]string

//Info contains commands and other goodies
var Info AbsoluteRoute

// This function will be called when the bot receives
// the "ready" event from Discord.
func ready(s *discordgo.Session, event *discordgo.Ready) {

	// Set the playing status.
	s.UpdateStatus(0, "!commands")
}

// This function will be called every time a new
// message is created on any channel that the autenticated bot has access to.
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.
	if m.Author.ID == s.State.User.ID {
		return
	}

	if strings.HasPrefix(m.Content, "!commands") {
		addCommands()
		var textBuild = ""
		for i := 0; i < len(Info.Commands); i++ {
			textBuild += Info.Commands[i].Command
			textBuild += "\n"
		}
		textBuild += "😂😂"
		s.ChannelMessageSend(m.ChannelID, string(len(Info.Commands))+textBuild)

	}

	if strings.Contains(m.Content, "😂") {
		s.ChannelMessageSend(m.ChannelID, "😂😂")
	}

	//Stops whatever sound that is playing on message origin server
	if strings.HasPrefix(m.Content, "!stop") {
		currentlyPlaying[m.GuildID] = ""
	}

	// check if the message is a command
	if strings.HasPrefix(m.Content, "!") && !strings.Contains(m.Content, "commands") && !strings.Contains(m.Content, "!stop") {

		if inUseServers[m.GuildID] != "" {
			log.Println("Server ", m.GuildID, " is possibly being command spammed")
			return
		}

		//Set the server as "Busy"
		inUseServers[m.GuildID] = m.GuildID

		var actualCommand = strings.Trim(m.Content, "!")

		insertCommandLog(actualCommand, m.Author.Username, m.GuildID)

		// Find the channel that the message came from.
		c, err := s.State.Channel(m.ChannelID)
		if err != nil {
			// Could not find channel.
			return
		}

		// Find the guild for that channel.
		g, err := s.State.Guild(c.GuildID)
		if err != nil {
			// Could not find guild.
			return
		}

		// Look for the message sender in that guild's current voice states.
		for _, vs := range g.VoiceStates {
			if vs.UserID == m.Author.ID {

				//Set as currently playing
				currentlyPlaying[m.GuildID] = m.GuildID
				err = playSound(s, g.ID, vs.ChannelID, actualCommand)

				//Clean up server info
				inUseServers[m.GuildID] = ""
				currentlyPlaying[m.GuildID] = ""

				if err != nil {
					log.Println("Error playing sound:", err)
				}

				return
			}
		}

	}

}

// This function will be called every time a new
// guild is joined.
func guildCreate(s *discordgo.Session, event *discordgo.GuildCreate) {

	//Server is down or unavailable
	if event.Guild.Unavailable {
		log.Println("Not able to connect to", event.Guild.ID)
		return
	}

	log.Println("Connected to", event.Guild.Name, event.Guild.ID)

}
