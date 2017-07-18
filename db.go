package main

import (
	"fmt"
	"log"

	"github.com/boltdb/bolt"
)

func updateDB(key, value string) error {
	db, err := bolt.Open(home()+"/mafia.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}

	db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("Mafia"))
		err := b.Put([]byte(key), []byte(value))
		return err
	})

	defer db.Close()
	return nil
}

func initDB() error {
	db, err := bolt.Open(home()+"/mafia.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}

	db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("Mafia"))
		if err != nil {
			return err
		}
		return nil
	})

	defer db.Close()
	return nil
}

func viewDB(key string) (string, error) {
	value := ""
	db, err := bolt.Open(home()+"/mafia.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}

	db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("Mafia"))
		value = fmt.Sprintf("%s", b.Get([]byte(key)))
		return nil
	})

	defer db.Close()
	return value, nil
}
