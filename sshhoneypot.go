package main

import (
	"encoding/base64"
	"fmt"
	"log"

	"database/sql"

	"github.com/gliderlabs/ssh"
	_ "github.com/mattn/go-sqlite3"
)

// DummyPasswordHandler is a password Handler that always rejects
// login attempts but will print a record

var db *sql.DB

func encodeString(s string) string {
	return base64.StdEncoding.EncodeToString([]byte(s))
}

func DummyPasswordHandler(c ssh.Context, pass string) bool {
	log.Printf("[TRAP] from=%s user=%s password=%s version=%s local=%s",
		c.RemoteAddr(),
		c.User(),
		pass,
		c.ClientVersion(),
		c.LocalAddr().String(),
	)
	// Encoding in base64
	add_record(c.RemoteAddr().String(), encodeString(c.User()), encodeString(pass), encodeString(c.ClientVersion()))
	return false
}

func add_record(sourceIP string, user string, password string, clientVersion string) {
	log.Println("Recording data: ", sourceIP, user, password, clientVersion)
	insertStmt := fmt.Sprintf(`INSERT INTO sshconnections VALUES (
		datetime('now'), "%s", "%s", "%s", "%s"
		);`, sourceIP, user, password, clientVersion)
	_, err := db.Exec(insertStmt)
	if err != nil {
		log.Fatalf("Failed to record transation: \nStatement: %s\nerror: %s", insertStmt, err)
	}
}

func initDB(dbname string) {
	var err error
	// Opening the DB
	db, err = sql.Open("sqlite3", dbname)
	if err != nil {
		log.Fatalf("Error while opening DB: %s with error: %s", dbname, err)
	}
	createStmt := `CREATE table IF NOT EXISTS sshconnections 
		(date text, source text, user text, password text, client text)
	`
	_, err = db.Exec(createStmt)
	if err != nil {
		log.Fatalf("error: %s", err)
	}
}

func main() {
	log.SetFlags(log.LstdFlags | log.LUTC)
	initDB("./mydb.sqlite")
	defer func() {
		log.Println("Closing DB and exiting...")
		db.Close()
	}()
	sshServer := &ssh.Server{
		Addr:            ":2222",
		Version:         "OpenSSH_8.1",
		PasswordHandler: DummyPasswordHandler,
	}
	log.Println("Starting SSH Honeypot.")
	log.Fatal(sshServer.ListenAndServe())

}
