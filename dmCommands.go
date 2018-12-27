package main

import (
	"database/sql"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/newtpuandre/dca"
)

type commandUpload struct {
	Command  string
	URL      string
	AuthorID string
}

//var commandUploadList []commandUpload
var commandUploadList map[string]*commandUpload

//func uploadAudio
//This ^ function will take a file by DM and add it to the commands

func uploadAudioInfo(s *discordgo.Session, m *discordgo.MessageCreate) {
	channel, err := s.UserChannelCreate(m.Author.ID)
	if err != nil {
		log.Println("Something went wrong whilst trying to create a DM, " + err.Error())
	}

	_, err = s.ChannelMessageSend(channel.ID, "Please upload the file to this chat. When it is confirmed uploaded please use the !commandname <your_name_here> to give it a command name. Only upload one sound at a time!")
	if err != nil {
		log.Println("Could not send message to discord in dmCommands, " + err.Error())
	}
}

func attachments(s *discordgo.Session, m *discordgo.MessageCreate) {
	channel, err := s.UserChannelCreate(m.Author.ID)
	if err != nil {
		log.Println("Something went wrong whilst trying to create a DM, " + err.Error())
	}

	if !strings.HasSuffix(m.Attachments[0].Filename, ".mp3") {
		_, err = s.ChannelMessageSend(channel.ID, "Only MP3 files are supported. Please convert it first")
		return
	}

	_, err = s.ChannelMessageSend(channel.ID, "File is sucessfully uploaded. Please use !commandname <your_name_here to give it a name! Example : !commandname jonas")
	if err != nil {
		log.Println("Could not send message to discord  " + err.Error())
	}

	var temp commandUpload
	temp.URL = m.Attachments[0].URL
	temp.AuthorID = m.Author.ID
	commandUploadList[m.Author.ID] = &temp

}

func commandName(s *discordgo.Session, m *discordgo.MessageCreate) {
	channel, err := s.UserChannelCreate(m.Author.ID)
	if err != nil {
		log.Println("Something went wrong whilst trying to create a DM, " + err.Error())
	}

	words := strings.Split(m.Content, " ")
	if len(words) > 1 {
		command := words[1]

		if commandUploadList[m.Author.ID].AuthorID != m.Author.ID {
			_, err = s.ChannelMessageSend(channel.ID, "Please upload a file first!")
			return
		}

		var index = -1
		for i := range DBSoundList {
			if DBSoundList[i].command == command {
				index = i
			}
		}

		if index != -1 {
			_, err = s.ChannelMessageSend(channel.ID, "Command is already in use")
			return
		}

		commandUploadList[m.Author.ID].Command = command

		err = downloadFile("./sounds/"+commandUploadList[m.Author.ID].Command+".mp3", commandUploadList[m.Author.ID].URL)
		if err != nil {
			log.Println(err)
		}

		// Encoding a file and saving it to disk
		var MP3Path = "./sounds/" + commandUploadList[m.Author.ID].Command + ".mp3"

		encodeSession, err := dca.EncodeFile(MP3Path, dca.StdEncodeOptions)
		// Make sure everything is cleaned up, that for example the encoding process if any issues happened isnt lingering around
		encodeSession.Cleanup()

		var fullpath = "./sounds/" + commandUploadList[m.Author.ID].Command + ".dca"

		output, err := os.Create(fullpath)
		if err != nil {
			log.Println(err)
		}

		io.Copy(output, encodeSession)

		os.Remove(MP3Path)

		//Insert file into DB.
		b, err := ioutil.ReadFile(fullpath) // just pass the file name
		if err != nil {
			log.Print(err)
			return
		}

		db, err := sql.Open("mysql", connectionString)
		if err != nil {
			log.Println(err)
			return
		}
		defer db.Close()

		stmtIns, err := db.Prepare("INSERT INTO sounds(filename,command,file) VALUES( ?, ? ,?)") // ? = placeholder
		if err != nil {
			log.Println(err)
			return
		}
		defer stmtIns.Close() // Close the statement when we leave main() / the program terminates

		_, err = stmtIns.Exec(fullpath, command, b) // Insert tuples (i, i^2)
		if err != nil {
			log.Println(err)
			return
		}

		loadFromList()

		_, err = s.ChannelMessageSend(channel.ID, "Please PM TPU to enable the sound! :joy: :joy:")
		if err != nil {
			log.Println("Could not send message to discord " + err.Error())
		}

	} else {
		_, err = s.ChannelMessageSend(channel.ID, "Something went wrong. Please try again later")
		if err != nil {
			log.Println("Could not send message to discord " + err.Error())
		}
	}
}

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

func downloadFile(filepath string, url string) (err error) {

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Check server response
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	// Writer the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}
