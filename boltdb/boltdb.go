package boltdb

import (
	bolt "github.com/boltdb/bolt"
	"log"
)

type DBService struct {
	DbRef *bolt.DB
}

func InitializeDB() *DBService {

	db, err := bolt.Open("telebot.db", 0600, nil)

	db.Update(func(tx *bolt.Tx) error {

		_, err := tx.CreateBucketIfNotExists([]byte("store"))

		if err != nil {
			log.Fatal(err)
			return err
		}
		return nil
	})

	if err != nil {
		log.Fatal(err)
		return nil
	}

	return &DBService{DbRef: db}
}

func (ds *DBService) InsertValue(key, value string) {

	ds.DbRef.Update(func(tx *bolt.Tx) error {

		b, err := tx.CreateBucketIfNotExists([]byte("store"))

		if err != nil {
			return err
		}
		b.Put([]byte(key), []byte(value))
		return nil
	})

}

func (ds *DBService) GetValue(key string) string {

	var value string

	ds.DbRef.Update(func(tx *bolt.Tx) error {

		b := tx.Bucket([]byte("store"))

		value = string(b.Get([]byte(key)))
		return nil
	})

	return value
}

func (ds *DBService) CloseConnection() error {

	err := ds.DbRef.Close()

	if err != nil {
		return err
	}
	return nil

}
