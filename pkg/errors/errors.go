package errors

import "fmt"

type AuthFailError struct {
	Username string
}

type MalformedToError struct {
	To string
}
type SendMailErorr struct{}

func (af *AuthFailError) Error() string {
	return fmt.Sprintf("Auth failed for user %s", af.Username)
}

func (mft *MalformedToError) Error() string {
	return fmt.Sprintf("Mailformed TO address: %s", mft.To)
}

func (sm *SendMailErorr) Error() string {
	return "Unable to send email via any of the MX records and ports"
}
