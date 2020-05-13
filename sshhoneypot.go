package main

import (
	"log"

	"github.com/gliderlabs/ssh"
)

// DummyPasswordHandler is a password Handler that always rejects
// login attempts but will print a record
func DummyPasswordHandler(c ssh.Context, pass string) bool {
	log.Printf("[TRAP] from=%s user=%s password=%s version=%s local=%s",
		c.RemoteAddr(),
		c.User(),
		pass,
		c.ClientVersion(),
		c.LocalAddr().String(),
	)
	return false
}

func main() {
	log.SetFlags(log.LstdFlags | log.LUTC)
	sshServer := &ssh.Server{
		Addr:            ":2222",
		Version:         "OpenSSH_8.1",
		PasswordHandler: DummyPasswordHandler,
	}
	log.Println("Starting SSH Honeypot.")
	log.Fatal(sshServer.ListenAndServe())

}
