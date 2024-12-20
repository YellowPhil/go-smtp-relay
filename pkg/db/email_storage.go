package db

import (
	"github.com/dgraph-io/badger/v4"
	"github.com/yellowphil/go-smtp-relay/pkg/utils"
)

type EmailStorage struct {
	*BasicStorage
}

func (es *EmailStorage) AllowedRcpt(email string) bool {
	err := es.db.View(func(txn *badger.Txn) error {
		_, err := txn.Get(utils.UnsafeBytes(email))
		return err

	})
	return (err == badger.ErrKeyNotFound) || (err == nil)
}

func (es *EmailStorage) AddAllowedRcpt(email string) error {
	// TODO: find a better way
	return es.Set(email, []byte(""))
}
