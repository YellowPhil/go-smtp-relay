package session

import (
	"io"
	"log"

	"github.com/emersion/go-smtp"
	errors "github.com/yellowphil/go-smtp-relay/pkg/errors"
)

const testUser = "testUser"
const testPassword = "testPassword"
const workerCount = 3

type Session struct {
	To       []string
	From     string
	Contents []byte
}

func (s *Session) Mail(from string, opts *smtp.MailOptions) error {
	log.Println("Mail from:", from)
	s.From = from
	return nil
}
func (s *Session) Rcpt(to string, options *smtp.RcptOptions) error {
	log.Println("Rcpt to: ", to)
	s.To = append(s.To, to)
	return nil
}

func (s *Session) Data(r io.Reader) error {
	if buffer, err := io.ReadAll(r); err != nil {
		return err
	} else {
		s.Contents = buffer
	}
	sendChan := make(chan *SendTask, len(s.To))
	workers := make([]*Worker, workerCount)
	for i := range workerCount {
		workers[i] = NewWorker(sendChan)
		go workers[i].SendWithRetries()
	}
	for _, to := range s.To {
		sendChan <- &SendTask{To: to, From: s.From, Contents: s.Contents}
	}
	close(sendChan)

	return nil
}

func (s *Session) Auth(username, password string) error {
	// TODO: change before release
	if username != testUser && password != testPassword {
		return &errors.AuthFailError{Username: username}
	}
	return nil
}

func (s *Session) Reset() {
	s.From = ""
	s.To = []string{}
	s.Contents = []byte{}
}

func (s *Session) Logout() error {
	return nil
}
