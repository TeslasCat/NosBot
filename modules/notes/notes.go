package notes

import (
	"log"
	"fmt"
	// "time"
	"github.com/lrstanley/girc"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

type Note struct {
	ID int
	Nick string
	Timestamp string
	Note string
}

func New (c *girc.Client, e girc.Event) {
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
		return
	}

	n := Note{1, e.Source.Name, e.Timestamp.String(), e.Trailing}

	sqlStmt = fmt.Sprintf("INSERT INTO notes(rowid, Nick, Timestamp, Note) VALUES(null,'%s','%s','%s');", n.Nick, n.Timestamp, n.Note)
	_, err = db.Exec(sqlStmt)
	if err != nil {
		log.Printf("Notes Module: %q: %s\n", err, sqlStmt)
		defer db.Close()
		return
	}

	c.Cmd.ReplyTo(e, girc.Fmt(fmt.Sprintf("ID: %d by %s on %s Note: %s", n.ID, n.Nick, n.Timestamp, n.Note)))
	return
}

func List (c *girc.Client, e girc.Event) {
	db, err := sql.Open("sqlite3", "./notes.db")
	if err != nil {
		log.Printf("Notes Module: ", err)
	}

	rows, err := db.Query("SELECT * FROM notes;")
	if err != nil {
		log.Printf("Notes Module: ", err)
		c.Cmd.ReplyTo(e, girc.Fmt("{red}No Notes"))
		return
	}
	defer rows.Close()
	for rows.Next() {
		var n Note
		err = rows.Scan(&n.ID, &n.Nick, &n.Timestamp, &n.Note)
		if err != nil {
			log.Printf("Notes Module: ", err)
		}
		c.Cmd.ReplyTo(e, girc.Fmt(fmt.Sprintf("ID: %d by %s on %s Note: %s", n.ID, n.Nick, n.Timestamp, n.Note)))

	}
	err = rows.Err()
	if err != nil {
		log.Printf("Notes Module: ", err)
	}
}