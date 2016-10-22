package db

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"

	pb "gopkg.in/cheggaaa/pb.v1"

	"github.com/boltdb/bolt"
)

//NewBucketNotFoundError formatting for missing bucket
func NewBucketNotFoundError(bucketName string) error {
	msg := fmt.Sprintf("Bucket not found: %s", bucketName)
	return errors.New(msg)
}

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

//CreateBuckets Creates buckets if they don't exist
func CreateBuckets(db *bolt.DB, bucketNames []string) {
	err := db.Update(func(tx *bolt.Tx) error {
		for _, bucketName := range bucketNames {
			_, err := tx.CreateBucketIfNotExists([]byte(bucketName))
			if err != nil {
				return NewCreateBucketError(bucketName, err)
			}
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
}

//UpdateStringList Puts a map[string][]string into the db
func UpdateStringList(db *bolt.DB, bucketName string, data map[string][]string) error {
	log.Printf("Importing %s", bucketName)
	bar := pb.New(len(data)).SetWidth(80).Start()
	bar.ShowSpeed = true

	err := db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(bucketName))
		if bucket == nil {
			return NewBucketNotFoundError(bucketName)
		}

		for key, listValues := range data {
			listValuesBytes, err := json.Marshal(listValues)
			if err != nil {
				return err
			}

			bucket.Put([]byte(key), listValuesBytes)
			bar.Increment()
		}

		return nil
	})

	bar.Finish()

	if err != nil {
		log.Fatal(err)
	}

	return nil
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
func GetAllKeys(db *bolt.DB, bucketName string) []string {
	var result []string
	err := db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(bucketName))
		if bucket == nil {
			return NewBucketNotFoundError(bucketName)
		}

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
func GetAll(db *bolt.DB, bucketName string) map[string]string {
	result := make(map[string]string)
	err := db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(bucketName))
		if bucket == nil {
			return NewBucketNotFoundError(bucketName)
		}

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
func GetAllLists(db *bolt.DB, bucketName string) map[string][]string {
	result := make(map[string][]string)
	err := db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(bucketName))
		if bucket == nil {
			return NewBucketNotFoundError(bucketName)
		}

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

//UniqueStrings convert map[string][]string in a set
func UniqueStrings(data map[string][]string) map[string]bool {
	result := make(map[string]bool)
	for _, elements := range data {
		for _, element := range elements {
			result[element] = true
		}
	}
	return result
}
