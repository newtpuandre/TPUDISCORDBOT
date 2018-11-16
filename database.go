package main

import (
	"database/sql"
	"io/ioutil"
	"log"
	"time"
)

func loadFromDB() {
	db, err := sql.Open("mysql", connectionString)
	if err != nil {
		panic(err.Error()) // Just for example purpose. You should use proper error handling instead of panic
	}
	defer db.Close()

	// Query the square-number of 13
	// Execute the query
	rows, err := db.Query("SELECT * FROM sounds")
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	// Get column names
	columns, err := rows.Columns()
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}

	// Make a slice for the values
	values := make([]sql.RawBytes, len(columns))

	// rows.Scan wants '[]interface{}' as an argument, so we must copy the
	// references into such a slice
	// See http://code.google.com/p/go-wiki/wiki/InterfaceSlice for details
	scanArgs := make([]interface{}, len(values))
	for i := range values {
		scanArgs[i] = &values[i]
	}

	for rows.Next() {
		// get RawBytes from data
		err = rows.Scan(scanArgs...)
		if err != nil {
			panic(err.Error()) // proper error handling instead of panic in your app
		}
		var tempDBSound DBSound

		var value string
		for i, col := range values {
			if col == nil {
				value = "NULL"
			} else if i == 4 && col != nil {
				// write the whole body at once
				err = ioutil.WriteFile(tempDBSound.filepath, col, 0644)
				if err != nil {
					panic(err)
				}
			} else {
				value = string(col)
			}

			switch i {
			case 0:
				tempDBSound.id = value
			case 1:
				tempDBSound.filepath = value
			case 2:
				tempDBSound.command = value
			case 3:
				tempDBSound.enabled = value
			}
		}
		addToList(tempDBSound)
	}
	if err = rows.Err(); err != nil {
		//panic(err.Error()) // proper error handling instead of panic in your app
		log.Println(err.Error())
	}
}

func insertCommandLog(command string, user string, serverid string) {
	db, err := sql.Open("mysql", connectionString)
	if err != nil {
		log.Println(err)
		return
	}
	defer db.Close()

	stmtIns, err := db.Prepare("INSERT INTO commandlog(command, user, timestamp, serverid) VALUES( ?, ?, ?, ?)") // ? = placeholder
	if err != nil {
		log.Println(err)
		return
	}
	defer stmtIns.Close() // Close the statement when we leave main() / the program terminates

	_, err = stmtIns.Exec(command, user, time.Now().Unix(), serverid) // Insert tuples (i, i^2)
	if err != nil {
		log.Println(err)
		return
	}

}

func addCommands() {
	var newInfo AbsoluteRoute
	Info = newInfo

	var addSound CommandRoute
	addSound.Command = commandText
	Info.Commands = append(Info.Commands, addSound)

	addSound.Command = "!github"
	Info.Commands = append(Info.Commands, addSound)

	for i := range DBSoundList {
		var temp CommandRoute
		if DBSoundList[i].enabled == "1" {
			temp.Command = "!" + DBSoundList[i].command
			Info.Commands = append(Info.Commands, temp)
		}
	}

}
