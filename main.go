package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

func determineListenAddress() (string, error) { //Inorder to get the port heroku assigns us
	port := os.Getenv("PORT")
	if port == "" {
		return "", fmt.Errorf("$PORT not set")
	}
	return ":" + port, nil
}

func main() {
	//Setup router
	router := mux.NewRouter().StrictSlash(true)

	inUseServers = make(map[string]string)
	currentlyPlaying = make(map[string]string)

	addRoutes(router)

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

	// Open the websocket and begin listening.
	err = dg.Open()
	if err != nil {
		log.Println("Error opening Discord session: ", err)
	}

	// Wait here until CTRL-C or other term signal is received.
	log.Println("TPU Discord bot is now running.  Press CTRL-C to exit.")

	//spesify IP and port if you want to run anything besides heroku
	log.Fatal(http.ListenAndServe("0.0.0.0:8080", router))

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session.
	dg.Close()
}
