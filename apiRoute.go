package main

import (
	"crypto/md5"
	"database/sql"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/newtpuandre/dca"
)

func addRoutes(r *mux.Router) {
	//Make the router handle routes. Routes located in routes.go
	r.HandleFunc("/", InfoRoute).Methods("GET")
	r.HandleFunc("/upload", upload)
	r.HandleFunc("/update", LoadFromDBRoute)
}

//InfoRoute returns a list of commands
func InfoRoute(w http.ResponseWriter, r *http.Request) {
	crutime := time.Now().Unix()
	h := md5.New()
	io.WriteString(h, strconv.FormatInt(crutime, 10))
	token := fmt.Sprintf("%x", h.Sum(nil))

	t, _ := template.ParseFiles("index.gtpl")
	t.Execute(w, token)
}

//LoadFromDBRoute updates memory with updated/new files
func LoadFromDBRoute(w http.ResponseWriter, r *http.Request) {
	loadFromList()
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.NewEncoder(w).Encode("Updated information")
}

func upload(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		crutime := time.Now().Unix()
		h := md5.New()
		io.WriteString(h, strconv.FormatInt(crutime, 10))
		token := fmt.Sprintf("%x", h.Sum(nil))

		t, _ := template.ParseFiles("upload.gtpl")
		t.Execute(w, token)

	} else {

		r.ParseMultipartForm(32 << 20)
		file, handler, err := r.FormFile("uploadfile")
		commandValue := r.FormValue("command")

		if strings.Contains(commandValue, "!") {
			commandValue = strings.Trim(commandValue, "!")
		}

		var index int
		index = -1
		for i := range DBSoundList {
			if DBSoundList[i].command == commandValue {
				index = i
			}
		}

		if index == -1 {
			fmt.Fprintf(w, "Command is already in use. Choose something else")
			return
		}

		if err != nil || !strings.HasSuffix(handler.Filename, ".mp3") {
			fmt.Fprintf(w, "File exention is not supported... Try again with a MP3 extension")
			return
		}
		defer file.Close()

		f, err := os.OpenFile("./sounds/"+handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			log.Println(err)
			return
		}

		// Encoding a file and saving it to disk
		var MP3Path = f.Name()
		defer f.Close()
		io.Copy(f, file)
		encodeSession, err := dca.EncodeFile(MP3Path, dca.StdEncodeOptions)
		// Make sure everything is cleaned up, that for example the encoding process if any issues happened isnt lingering around
		defer encodeSession.Cleanup()

		var withoutMP3 = strings.Split(f.Name(), ".mp3")
		var fullpath = withoutMP3[0] + ".dca"

		output, err := os.Create(fullpath)
		if err != nil {
			log.Println(err)
		}

		io.Copy(output, encodeSession)

		os.Remove(MP3Path)

		fmt.Fprintf(w, "sound was added to the bot with the command !%s", commandValue)

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

		_, err = stmtIns.Exec(fullpath, commandValue, b) // Insert tuples (i, i^2)
		if err != nil {
			log.Println(err)
			return
		}

		loadSound(fullpath, commandValue)
	}
}
