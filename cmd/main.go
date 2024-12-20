package main

import (
	"log"
	"time"

	"github.com/emersion/go-smtp"
	"github.com/yellowphil/go-smtp-relay/pkg/config"
	"github.com/yellowphil/go-smtp-relay/pkg/session"
)

type Backend struct {
	cfg *config.Config
}

func (b *Backend) NewSession(_ *smtp.Conn) (smtp.Session, error) {
	return &session.Session{
		Cfg: *b.cfg,
	}, nil
}

func SMTPServer(config *config.Config) *smtp.Server {
	server := smtp.NewServer(&Backend{cfg: config})
	server.Domain = config.Connection.ListenDomain
	server.Addr = config.Connection.ListenAddr
	server.WriteTimeout = 10 * time.Second
	server.ReadTimeout = 10 * time.Second

	server.AllowInsecureAuth = config.Connection.AllowInsecureAuth
	server.MaxMessageBytes = 1024 * 1024
	server.MaxRecipients = 20

	return server

}

func main() {
	cfg, err := config.NewConfigFromFile("./config.yaml")
	if err != nil {
		log.Fatal(err)
	}
	server := SMTPServer(cfg)

	log.Println("Starting server at", server.Addr)
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
