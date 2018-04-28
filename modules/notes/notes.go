package notes

import (
	"../../types"
	"log"
	"fmt"
	"strconv"
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
			response = list(message.Arguments)
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
	sqlStmt := `CREATE TABLE IF NOT EXISTS notes (ID INTEGER NOT NULL PRIMARY KEY, Nick TEXT, Timestamp DATETIME DEFAULT CURRENT_TIMESTAMP, Note TEXT);`
	_, err = db.Exec(sqlStmt)
	if err != nil {
		log.Printf("Notes Module: %q: %s\n", err, sqlStmt)
		defer db.Close()
		return response
	}

	n := Note{1, message.Nick, message.Timestamp, message.Message}

	sqlStmt = fmt.Sprintf("INSERT INTO notes(rowid, Nick, Note) VALUES(null,'%s','%s');", n.Nick, n.Note)
	_, err = db.Exec(sqlStmt)
	if err != nil {
		log.Printf("Notes Module: %q: %s\n", err, sqlStmt)
		defer db.Close()
		return response
	}

	response.Type = "action"
	response.Message = "adds to list"
	return response
}

func list (arguments []string) types.Response {
	response := types.Response{}

	db, err := sql.Open("sqlite3", "./notes.db")
	if err != nil {
		log.Printf("Notes Module: ", err)
	}

	sql := "SELECT ID, Nick, strftime('%d/%m/%Y %H:%M', Timestamp) AS Timestamp, Note FROM notes"

	if len(arguments) > 0 {
		limit, err := strconv.Atoi(arguments[0]);
		if (err == nil) {
			sql += fmt.Sprintf(" ORDER BY Timestamp DESC LIMIT %d ", limit)
		} else {
			// Convert argument to date range
			switch {
				case arguments[0] == "today":
					sql += " WHERE Timestamp > datetime('now', 'start of day')"
				case arguments[0] == "yesterday":
					sql += " WHERE Timestamp BETWEEN datetime('now', '-1 days') AND datetime('now', 'start of day')"
			}

		}
	}

	rows, err := db.Query(sql)
	if err != nil {
		log.Printf("Notes Module: ", err)
		response.Message = "{red}No Notes"
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
