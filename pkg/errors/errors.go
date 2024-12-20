package errors

import "fmt"

type AuthFailError struct {
	Username string
}

type MalformedToError struct {
	To string
}
type SendMailError struct{}

type NoConnectionError struct{}

type MXLookupFailError struct{}

func (af *AuthFailError) Error() string {
	return fmt.Sprintf("auth failed for user %s", af.Username)
}

func (mft *MalformedToError) Error() string {
	return fmt.Sprintf("mailformed TO address: %s", mft.To)
}

func (sm *SendMailError) Error() string {
	return "unable to send email via any of the MX records and ports"
}

func (nc *NoConnectionError) Error() string {
	return "connection type was not chosen before sending email"
}

func (mxl *MXLookupFailError) Error() string {
	return "unable to perform MX lookup on receivers domain"
}
