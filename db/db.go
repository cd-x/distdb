package db

import (
	"github.com/boltdb/bolt"
)

var defaultBucketName = []byte("default")

type Database struct {
	db *bolt.DB
}

func NewDatabase(path string) (db *Database, close func() error, err error) {
	boltdb, err := bolt.Open(path, 0600, nil)
	if err != nil {
		return nil, nil, err
	}
	closeFunc := boltdb.Close
	return &Database{db: boltdb}, closeFunc, nil
}

func (d *Database) SetKey(key string, value []byte) error {
	return d.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(defaultBucketName)
		err := b.Put([]byte(key), []byte(value))
		return err
	})
}

func (d *Database) GetKey(key string) ([]byte, error) {
	var value []byte
	err := d.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(defaultBucketName)
		value = b.Get([]byte(key))
		return nil
	})
	if err != nil {
		return nil, err
	} else {
		return value, nil
	}
}
