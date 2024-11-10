package main

import (
	"log"
	"time"

	"github.com/emersion/go-smtp"
	"github.com/yellowphil/go-smtp-relay/pkg/config"
	"github.com/yellowphil/go-smtp-relay/pkg/session"
)

type Backend struct {
	cfg config.Config
}

func (b *Backend) NewSession(_ *smtp.Conn) (smtp.Session, error) {
	return &session.Session{
		Cfg: b.cfg,
	}, nil
}

func SMTPServer(listenAddr string, domain string, debug bool) *smtp.Server {
	server := smtp.NewServer(&Backend{})
	server.Domain = domain
	server.Addr = listenAddr
	server.WriteTimeout = 10 * time.Second
	server.ReadTimeout = 10 * time.Second

	server.AllowInsecureAuth = debug
	server.MaxMessageBytes = 1024 * 1024
	server.MaxRecipients = 20

	return server

}

func main() {
	server := SMTPServer("127.0.0.1:2525", "localhost", true)

	log.Println("Starting server at", server.Addr)
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
