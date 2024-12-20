package db

import (
	"sync"
	"time"

	"github.com/dgraph-io/badger/v4"
	"github.com/yellowphil/go-smtp-relay/pkg/utils"
)

var userStorage *UserStorage = nil
var emailStorage *EmailStorage = nil

var userLock *sync.Once
var emailLock *sync.Once

type BasicStorage struct {
	db         *badger.DB
	gcInterval time.Duration
	done       chan struct{}
}

func (s *BasicStorage) Get(key string) (value []byte, err error) {
	if len(key) <= 0 {
		return nil, nil
	}
	err = s.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get(utils.UnsafeBytes(key))
		if err != nil {
			if err == badger.ErrKeyNotFound {
				return nil
			}
			return err
		}
		value, err = item.ValueCopy(nil)
		return err
	})
	return
}

func (s *BasicStorage) Set(key string, val []byte) (err error) {
	if len(key) <= 0 || len(val) <= 0 {
		return nil
	}

	err = s.db.Update(func(txn *badger.Txn) error {
		e := badger.NewEntry(utils.UnsafeBytes(key), val)
		return txn.SetEntry(e)
	})
	return
}

func (s *BasicStorage) Delete(key string) (err error) {
	if len(key) <= 0 {
		return nil
	}
	err = s.db.Update(func(txn *badger.Txn) error {
		return txn.Delete(utils.UnsafeBytes(key))
	})
	return
}

func (s *BasicStorage) Close() error {
	s.done <- struct{}{}
	return s.db.Close()
}

func (s *BasicStorage) gcLoop() {
	ticker := time.NewTicker(s.gcInterval)
	defer ticker.Stop()
	for {
		select {
		case <-s.done:
			return
		case <-ticker.C:
			_ = s.db.RunValueLogGC(0.7)
		}
	}
}

func GetUserStorage() *UserStorage {
	userLock.Do(func() {
		userStorage = &UserStorage{}
	})
	return userStorage
}

func GetEmailStorage() *EmailStorage {
	emailLock.Do(func() {
		emailStorage = &EmailStorage{}
	})
	return emailStorage
}
