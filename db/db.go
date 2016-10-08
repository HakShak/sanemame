package db

import (
	"github.com/boltdb/bolt"
	"log"
)

func GetAllKeys(db *bolt.DB, bucket string) []string {
	var result []string
	err := db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(bucket))

		bucket.ForEach(func(k, v []byte) error {
			result = append(result, string(k))
			return nil
		})

		return nil
	})
	if err != nil {
		log.Fatal(err)
	}

	return result
}

func GetAll(db *bolt.DB, bucket string) map[string]string {
	result := make(map[string]string)
	err := db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(bucket))

		bucket.ForEach(func(k, v []byte) error {
			result[string(k)] = string(v)
			return nil
		})

		return nil
	})
	if err != nil {
		log.Fatal(err)
	}

	return result
}
