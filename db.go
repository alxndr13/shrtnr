package main

import (
	"fmt"

	bolt "go.etcd.io/bbolt"
)

func (a *App) saveToDatabase(shortCode string, url string) error {
	err := a.Db.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte("shortCodes"))
		if err != nil {
			return err
		}
		bucket.Put([]byte(shortCode), []byte(url))
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func (a *App) getFromDatabase(shortCode string) (url string, err error) {
	err = a.Db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("shortCodes"))
		if bucket == nil {
			return err
		}
		result := bucket.Get([]byte(shortCode))
		if len(result) == 0 {
			return fmt.Errorf("no result in database")
		}
		fmt.Println("Result:", string(result))
		url = string(result)
		return nil
	})

	if err != nil {
		return "", err
	}
	return url, nil

}

func (a *App) createDatabaseIfNotExists() error {
	db, err := bolt.Open(a.DbPath, 0600, nil)
	if err != nil {
		return err
	}
	a.Db = db
	return nil
}

func (a *App) getAmountOfLinks() (int, error) {
	var amount int
	err := a.Db.View(func(tx *bolt.Tx) error {
		// Assume bucket exists and has keys
		b := tx.Bucket([]byte("shortCodes"))

		c := b.Cursor()

		for k, _ := c.First(); k != nil; k, _ = c.Next() {
			amount++
		}

		return nil
	})
	return amount, err
}
