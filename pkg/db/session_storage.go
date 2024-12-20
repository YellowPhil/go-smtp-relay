package db

import (
	"github.com/dgraph-io/badger/v4"
	"github.com/yellowphil/go-smtp-relay/pkg/utils"
	"time"
)

type BadgerSessionStorage struct {
	*BasicStorage
}

func (s *BadgerSessionStorage) Set(key string, val []byte, exp time.Duration) (err error) {
	if len(key) <= 0 || len(val) <= 0 {
		return nil
	}

	err = s.db.Update(func(txn *badger.Txn) error {
		e := badger.NewEntry(utils.UnsafeBytes(key), val).WithTTL(exp)
		return txn.SetEntry(e)
	})
	return
}

func (s *BadgerSessionStorage) Reset() error {
	return s.db.DropAll()
}
