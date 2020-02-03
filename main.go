package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	_ "github.com/go-sql-driver/mysql"
)

func main() {

	inUseServers = make(map[string]string)
	currentlyPlaying = make(map[string]string)
	commandUploadList = make(map[string]*commandUpload)

	loadFromList()

	// Create a new Discord session using the provided bot token.
	dg, err := discordgo.New("Bot " + DiscToken)
	if err != nil {
		log.Println("Error creating Discord session: ", err)
		return
	}

	// Register ready as a callback for the ready events.
	dg.AddHandler(ready)

	// Register messageCreate as a callback for the messageCreate events.
	dg.AddHandler(messageCreate)

	// Register guildCreate as a callback for the guildCreate events.
	dg.AddHandler(guildCreate)

	err = dg.Open()
	if err != nil {
		log.Println("Error opening Discord session: ", err)
	}

	// Wait here until CTRL-C or other term signal is received.
	log.Println("TPU Discord bot is now running.  Press CTRL-C to exit.")

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session.
	dg.Close()
}

func cleanup(GuildID string) {
	inUseServers[GuildID] = ""
	currentlyPlaying[GuildID] = ""
	log.Println("Setting server as AVAILABLE:", GuildID)
}
