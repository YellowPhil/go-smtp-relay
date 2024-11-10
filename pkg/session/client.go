package session

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/smtp"
	"strings"

	"github.com/yellowphil/go-smtp-relay/pkg/errors"
)

type Client struct{}

func MXLookup(domain string) ([]*net.MX, error) {
	mxRecords, err := net.LookupMX(domain)
	if err != nil {
		return nil, fmt.Errorf("error looking up MX for domain %s", domain)
	}
	return mxRecords, nil
}

func (c *Client) sendMail(client *smtp.Client, from, to string, data []byte) error {
	if err := client.Mail(from); err != nil {
		return err
	}
	if err := client.Rcpt(to); err != nil {
		return err
	}
	writer, err := client.Data()
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
	toParts := strings.Split(to, "@")
	if len(toParts) != 2 {
		return &errors.MalformedToError{To: to}
	}
	domain := toParts[1]
	mxRecords, err := MXLookup(domain)
	if err != nil {
		return err
	}

	for _, mx := range mxRecords {
		host := mx.Host
		// use this sequence to go from most to less secure
		for _, port := range []int{587, 465, 25} {
			addr := fmt.Sprintf("%s:%d", host, port)
			var client *smtp.Client
			var err error

			switch port {
			case 587:
				// SMTPS
				tlsConfig := &tls.Config{ServerName: host}
				conn, err := tls.Dial("tcp", addr, tlsConfig)
				if err != nil {
					continue
				}
				client, err = smtp.NewClient(conn, host)
				if err != nil {
					continue
				}

			case 25, 465:
				client, err = smtp.Dial(addr)
				if err != nil {
					continue
				}

				if port == 587 {
					if err = client.StartTLS(&tls.Config{ServerName: host}); err != nil {
						client.Close()
						continue
					}
				}
			}
			if err != nil {
				continue
			}
			if err := c.sendMail(client, to, from, data); err != nil {
				continue
			}
			client.Quit()
			return nil
		}
	}
	return &errors.SendMailErorr{}
}

func (c *Client) Send(to, from string, data []byte, retries int) {
	go func() {
		for range retries {
			if err := c.SendMail(from, to, data); err != nil {
				return
			}
		}
	}()
}
