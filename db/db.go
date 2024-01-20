package db

import (
	"fmt"

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
	db = &Database{db: boltdb}
	closeFunc := boltdb.Close

	if err := db.createBucket(); err != nil {
		closeFunc()
		return nil, nil, fmt.Errorf("create bucket %s failed: %v", path, err)
	}
	return &Database{db: boltdb}, closeFunc, nil
}

// Methos to create bucket in database if doesn't exist
func (d *Database) createBucket() error {
	return d.db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(defaultBucketName)
		if err != nil {
			fmt.Printf("create bucket failed at create bucket %s: %v", defaultBucketName, err)
		}
		return err
	})
}

// API: SetKey puts the key in the database
func (d *Database) SetKey(key string, value []byte) error {
	return d.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(defaultBucketName)
		err := b.Put([]byte(key), []byte(value))
		return err
	})
}

// API: GetKey returns value from database
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
