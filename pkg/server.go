package server

import (
	"log"
	"time"

	smtp "github.com/emersion/go-smtp"
	session "github.com/yellowphil/go-smtp-relay/pkg/session"
)

type Backend struct{}

func (b *Backend) NewSession(_ *smtp.Conn) (smtp.Session, error) {
	return &session.Session{}, nil
}

func New() *smtp.Server {
	server := smtp.NewServer(&Backend{})
	server.Addr = ":2525"
	server.Domain = "localhost"
	server.WriteTimeout = 10 * time.Second
	server.ReadTimeout = 10 * time.Second
	server.MaxMessageBytes = 1024 * 1024
	server.MaxRecipients = 20
	//TODO: change for future release
	server.AllowInsecureAuth = true

	log.Println("Starting server at ", server.Addr)
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
	return server
}
