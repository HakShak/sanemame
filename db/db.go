package db

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/boltdb/bolt"
)

//NewPutBucketError formatting for create bucket errors
func NewPutBucketError(bucketName string, key string, value string, err error) error {
	msg := fmt.Sprintf("Putting %s in %s|%s failed: %s", value, bucketName, key, err)
	return errors.New(msg)
}

//NewCreateBucketError formatting for create bucket errors
func NewCreateBucketError(bucketName string, err error) error {
	msg := fmt.Sprintf("Create bucket failed for %s: %s", bucketName, err)
	return errors.New(msg)
}

//InList Test membership of list
func InList(list []string, key string) bool {
	for _, check := range list {
		if check == key {
			return true
		}
	}
	return false
}

//GetAllKeys Helper to get all keys in a bucket
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

//GetAll Helper to get all keys and values in a bucket
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

//GetAllLists Helper to get all lists in a bucket
func GetAllLists(db *bolt.DB, bucket string) map[string][]string {
	result := make(map[string][]string)
	err := db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(bucket))

		bucket.ForEach(func(k, v []byte) error {
			var newList []string
			err := json.Unmarshal(v, &newList)
			result[string(k)] = newList
			if err != nil {
				log.Fatal(err)
			}
			return nil
		})
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}

	return result
}
