package session

import (
	"io"
	"log"

	"github.com/emersion/go-smtp"
	"github.com/yellowphil/go-smtp-relay/pkg/config"
	"github.com/yellowphil/go-smtp-relay/pkg/errors"
)

type Session struct {
	To       []string
	From     string
	Contents []byte
	Cfg      config.Config
	client   Client
}

func (s *Session) Mail(from string, opts *smtp.MailOptions) error {
	log.Println("Mail from:", from)

	s.From = from
	return nil
}
func (s *Session) Rcpt(to string, options *smtp.RcptOptions) error {
	log.Println("Rcpt to: ", to)

	return nil
}

func (s *Session) Data(r io.Reader) error {
	if buffer, err := io.ReadAll(r); err != nil {
		return err
	} else {
		s.Contents = buffer
	}
	for _, to := range s.To {
		s.client.UseSMTPS()
		if err := s.client.SendMail(s.From, to, s.Contents); err == nil {
			continue
		}
		s.client.UseSTARTTLS()
		if err := s.client.SendMail(s.From, to, s.Contents); err == nil {
			continue
		}
		if s.Cfg.AllowInsecure {
			s.client.UseInsecure()
			if err := s.client.SendMail(s.From, to, s.Contents); err == nil {
				continue
			}
		}
		log.Printf("ERROR: could not send an email to %s\n", to)
	}
	return nil
}

func (s *Session) Auth(username, password string) error {
	if username != s.Cfg.Creds.Username && password != s.Cfg.Creds.Password {
		return &errors.AuthFailError{Username: username}
	}
	return nil
}

func (s *Session) Reset() {
	s.From = ""
	s.To = []string{}
	s.Contents = nil
}

func (s *Session) Logout() error {
	return nil
}
