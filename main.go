package main

import (
	"log"
	"time"

	"github.com/emersion/go-smtp"
	"github.com/yellowphil/go-smtp-relay/pkg/session"
)

type Backend struct {
}

func (b *Backend) NewSession(_ *smtp.Conn) (smtp.Session, error) {
	return &session.Session{}, nil
}

func main() {
	server := smtp.NewServer(&Backend{})
	server.Domain = "localhost"
	server.Addr = "127.0.0.1:2525"
	server.WriteTimeout = 10 * time.Second
	server.ReadTimeout = 10 * time.Second
	// TODO: fix before production
	server.AllowInsecureAuth = true
	server.MaxMessageBytes = 1024 * 1024
	server.MaxRecipients = 20

	log.Println("Starting server at", server.Addr)
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
