package db

import (
	"fmt"
	"github.com/yellowphil/go-smtp-relay/pkg/utils"
)

type UserStorage struct {
	*BasicStorage
}

func (db *UserStorage) AddUser(username, password string) error {
	if username, err := db.Get(username); err != nil || username != nil {
		return fmt.Errorf("username already exists")
	}
	pwdHash := utils.Sha3SumString(password)
	db.Set(username, pwdHash)

	return nil
}

func (db *UserStorage) GetUser(username string) ([]byte, error) {
	if password, err := db.Get(username); err != nil {
		return nil, err
	} else {
		if password == nil {
			return nil, fmt.Errorf("username not found")
		}
		return password, nil
	}
}
