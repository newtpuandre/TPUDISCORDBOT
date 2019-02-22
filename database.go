package main

import (
	"database/sql"
	"io/ioutil"
	"log"
)

func loadFromDB() {
	db, err := sql.Open("mysql", connectionString)
	if err != nil {
		log.Println(err.Error()) // Just for example purpose. You should use proper error handling instead of panic
	}
	defer db.Close()

	// Execute the query
	rows, err := db.Query("SELECT * FROM sounds")
	if err != nil {
		log.Println(err.Error())
	}
	// Get column names
	columns, err := rows.Columns()
	if err != nil {
		log.Println(err.Error())
	}

	// Make a slice for the values
	values := make([]sql.RawBytes, len(columns))

	scanArgs := make([]interface{}, len(values))
	for i := range values {
		scanArgs[i] = &values[i]
	}

	for rows.Next() {
		// get RawBytes from data
		err = rows.Scan(scanArgs...)
		if err != nil {
			log.Println(err.Error())
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
					log.Println(err)
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
		log.Println(err.Error())
	}
}

func DBenableCommand(command string) (err error) {
	db, err := sql.Open("mysql", connectionString)
	if err != nil {
		log.Println(err)
		return err
	}
	defer db.Close()

	stmtIns, err := db.Prepare("UPDATE sounds SET ENABLED = 1 WHERE COMMAND=?") // ? = placeholder
	if err != nil {
		log.Println(err)
		return err
	}
	defer stmtIns.Close() // Close the statement when we leave main() / the program terminates

	_, err = stmtIns.Exec(command)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func DBdisableCommand(command string) (err error) {
	db, err := sql.Open("mysql", connectionString)
	if err != nil {
		log.Println(err)
		return err
	}
	defer db.Close()

	stmtIns, err := db.Prepare("UPDATE sounds SET ENABLED = 0 WHERE COMMAND=?") // ? = placeholder
	if err != nil {
		log.Println(err)
		return err
	}
	defer stmtIns.Close() // Close the statement when we leave main() / the program terminates

	_, err = stmtIns.Exec(command)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
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
