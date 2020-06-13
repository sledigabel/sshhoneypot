package main

import (
	"encoding/base64"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"database/sql"

	"github.com/gliderlabs/ssh"
	_ "github.com/mattn/go-sqlite3"
)

// DummyPasswordHandler is a password Handler that always rejects
// login attempts but will print a record

var db *sql.DB
var last_id uint
var db_lock *sync.Mutex

type db_entry struct {
	RemoteAddr    string
	Username      string
	Password      string
	Clientversion string
	Localaddr     string
}

func encodeString(s string) string {
	return base64.StdEncoding.EncodeToString([]byte(s))
}

func decodeString(s string) string {
	r, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		log.Fatal("Decoding string failed", err)
	}
	return string(r)
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
	entry := &db_entry{
		RemoteAddr:    c.RemoteAddr().String(),
		Username:      encodeString(c.User()),
		Password:      encodeString(pass),
		Clientversion: encodeString(c.ClientVersion()),
		Localaddr:     c.LocalAddr().String(),
	}
	add_record(entry)
	return false
}

// func add_record(sourceIP string, user string, password string, clientVersion string) {
func add_record(record *db_entry) {
	db_lock.Lock()
	defer db_lock.Unlock()
	// log.Println("Recording data: ", record.RemoteAddr, record.Username, record.Password, record.Clientversion)
	insertStmt := fmt.Sprintf(`INSERT INTO sshconnections (date, source, user, password, client) VALUES (
		datetime('now'), "%s", "%s", "%s", "%s"
		);`, record.RemoteAddr, record.Username, record.Password, record.Clientversion)
	_, err := db.Exec(insertStmt)
	if err != nil {
		log.Fatalf("Failed to record transation: \nStatement: %s\nerror: %s", insertStmt, err)
	}
}

func list_records() ([]*db_entry, error) {
	db_lock.Lock()
	defer db_lock.Unlock()
	log.Printf("Querying data from id: %d", last_id)
	row, err := db.Query("SELECT id, source, user, password, client FROM sshconnections WHERE id > ?", last_id)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer row.Close()
	return_rows := make([]*db_entry, 0)
	var id uint
	for row.Next() { // Iterate and fetch the records from result cursor
		r := &db_entry{}
		row.Scan(&id, &r.RemoteAddr, &r.Username, &r.Password, &r.Clientversion)
		return_rows = append(return_rows, r)
		last_id = id
	}
	return return_rows, nil
}

func report() {
	records, err := list_records()
	if err != nil {
		log.Println("Error during the fetching of results", err)
	}
	for _, record := range records {
		record.Clientversion = decodeString(record.Clientversion)
		record.Username = decodeString(record.Username)
		record.Password = decodeString(record.Password)
		log.Println(record)
	}
}

func initDB(dbname string) {
	var err error
	last_id = 0
	db_lock = &sync.Mutex{}
	db_lock.Lock()
	defer db_lock.Unlock()
	// Opening the DB
	db, err = sql.Open("sqlite3", dbname)
	if err != nil {
		log.Fatalf("Error while opening DB: %s with error: %s", dbname, err)
	}
	createStmt := `CREATE table IF NOT EXISTS sshconnections
		(id integer primary key autoincrement,
		date text,
		source text,
		user text,
		password text,
		client text
		)
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

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)
	log.Println("Starting SSH Honeypot.")
	go func() { log.Fatal(sshServer.ListenAndServe()) }()
	select {
	case <-c:
		log.Println("Caught Interrupt signal. Terminating...")
		sshServer.Close()
		os.Exit(0)
	}

}
