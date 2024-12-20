package session

import (
	"crypto/tls"
	"fmt"
	"net/smtp"
)

type Connection interface {
	Connect(addr string) (*smtp.Client, error)
}

type SMTPSConnection struct{}

type STARTTLSConnection struct{}

type InsecureConnection struct{}

func (c *SMTPSConnection) Connect(addr string) (*smtp.Client, error) {
	const smtpsPort = 465
	host := fmt.Sprintf("%s:%s", addr, smtpsPort)

	tlsClient := &tls.Config{ServerName: host}
	conn, err := tls.Dial("tcp", addr, tlsClient)
	if err != nil {
		return nil, err
	}
	return smtp.NewClient(conn, host)
}

func (c *STARTTLSConnection) Connect(addr string) (*smtp.Client, error) {
	const starTLSPort = 465
	host := fmt.Sprintf("%s:%d", addr, starTLSPort)

	client, err := smtp.Dial(host)
	if err != nil {
		return nil, err
	}
	if err := client.StartTLS(&tls.Config{ServerName: host}); err != nil {
		return nil, err
	}
	return client, err
}

func (c *InsecureConnection) Connect(addr string) (*smtp.Client, error) {
	return smtp.Dial(addr)
}
