package notes

import (
	"../../types"
	"log"
	"fmt"
	// "time"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

type Note struct {
	ID int
	Nick string
	Timestamp string
	Note string
}

func Handle (message types.Message) types.Response {
	var response types.Response

	switch {
		case message.Command == "note":
			response = new(message)
		case message.Command == "list":
			response = list()
	}

	return response
}

func new (message types.Message) types.Response {
	response := types.Response{}

	db, err := sql.Open("sqlite3", "./notes.db")
	if err != nil {
		log.Printf("Notes Module: ", err)
	}
	defer db.Close()
	sqlStmt := `CREATE TABLE IF NOT EXISTS notes (ID INTEGER NOT NULL PRIMARY KEY, Nick TEXT, Timestamp TEXT, Note TEXT);`
	_, err = db.Exec(sqlStmt)
	if err != nil {
		log.Printf("Notes Module: %q: %s\n", err, sqlStmt)
		defer db.Close()
		return response
	}

	n := Note{1, message.Nick, message.Timestamp, message.Message}

	sqlStmt = fmt.Sprintf("INSERT INTO notes(rowid, Nick, Timestamp, Note) VALUES(null,'%s','%s','%s');", n.Nick, n.Timestamp, n.Note)
	_, err = db.Exec(sqlStmt)
	if err != nil {
		log.Printf("Notes Module: %q: %s\n", err, sqlStmt)
		defer db.Close()
		return response
	}

	response.Type = "action"
	response.Messages = []string{"adds to list"}
	return response
}

func list () types.Response {
	response := types.Response{}

	db, err := sql.Open("sqlite3", "./notes.db")
	if err != nil {
		log.Printf("Notes Module: ", err)
	}

	rows, err := db.Query("SELECT * FROM notes;")
	if err != nil {
		log.Printf("Notes Module: ", err)
		response.Messages = []string{"{red}No Notes"}
		return response
	}

	defer rows.Close()
	for rows.Next() {
		var n Note
		err = rows.Scan(&n.ID, &n.Nick, &n.Timestamp, &n.Note)
		if err != nil {
			log.Printf("Notes Module: ", err)
		}

		response.Messages = append(response.Messages, fmt.Sprintf("ID: %d by %s on %s Note: %s", n.ID, n.Nick, n.Timestamp, n.Note))
	}

	err = rows.Err()
	if err != nil {
		log.Printf("Notes Module: ", err)
	}
	
	return response
}
