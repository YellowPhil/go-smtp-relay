package session

import (
	"fmt"
	"github.com/yellowphil/go-smtp-relay/pkg/errors"
	"net"
	"net/smtp"
	"strings"
)

type Client struct {
	smtp       *smtp.Client
	connection Connection
}

func MXLookup(domain string) ([]*net.MX, error) {
	mxRecords, err := net.LookupMX(domain)
	if err != nil {
		return nil, fmt.Errorf("error looking up MX for domain %s", domain)
	}
	return mxRecords, nil
}

func (c *Client) UseSMTPS() {
	c.connection = &SMTPSConnection{}
}
func (c *Client) UseSTARTTLS() {
	c.connection = &STARTTLSConnection{}
}
func (c *Client) UseInsecure() {
	c.connection = &InsecureConnection{}
}

func (c *Client) sendMail(from, to string, data []byte) error {
	if err := c.smtp.Mail(from); err != nil {
		return err
	}
	if err := c.smtp.Rcpt(to); err != nil {
		return err
	}
	writer, err := c.smtp.Data()
	if err != nil {
		writer.Close()
		return err
	}
	if _, err := writer.Write(data); err != nil {
		writer.Close()
		return err
	}
	return writer.Close()
}

func (c *Client) SendMail(from, to string, data []byte) error {
	if c.connection == nil {
		return &errors.NoConnectionError{}
	}
	toParts := strings.Split(to, "@")
	if len(toParts) != 2 {
		return &errors.MalformedToError{To: to}
	}
	domain := toParts[1]
	mxRecords, err := MXLookup(domain)
	if err != nil {
		return &errors.MXLookupFailError{}
	}

	for _, mx := range mxRecords {
		var host = mx.Host
		var smtpClient *smtp.Client
		var err error

		if smtpClient, err = c.connection.Connect(host); err != nil {
			continue
		}
		c.smtp = smtpClient
		if err = c.sendMail(from, to, data); err == nil {
			return nil
		}
	}
	return &errors.SendMailErorr{}
}
