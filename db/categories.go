package db

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/HakShak/sanemame/filetypes/catverini"
	"github.com/boltdb/bolt"
	pb "gopkg.in/cheggaaa/pb.v1"
)

//CategoryRawMachines bucket name
const CategoryRawMachines = "categoryraw-machines"

//CategoryPrimaryMachines bucket name
const CategoryPrimaryMachines = "categoryprimary-machines"

//CategorySecondaryMachines bucket name
const CategorySecondaryMachines = "categorysecondary-machines"

func update(db *bolt.DB, bucketName string, data map[string][]string) error {
	log.Printf("Importing %s", bucketName)
	bar := pb.New(len(data)).SetWidth(80).Start()
	bar.ShowSpeed = true

	err := db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(bucketName))
		if bucket == nil {
			msg := fmt.Sprintf("Bucket not found: %s", bucketName)
			return errors.New(msg)
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

//UpdateCategories populates categories from INI
func UpdateCategories(db *bolt.DB, fileName string) {

	err := db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(CategoryRawMachines))
		if err != nil {
			return err
		}
		_, err = tx.CreateBucketIfNotExists([]byte(CategoryPrimaryMachines))
		if err != nil {
			return err
		}
		_, err = tx.CreateBucketIfNotExists([]byte(CategorySecondaryMachines))
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}

	rawCategories := make(map[string][]string)
	primaryCategories := make(map[string][]string)
	secondaryCategories := make(map[string][]string)

	catverini.Load(fileName, catverini.EntryRead(
		func(machine string, category *catverini.Category) error {
			rawCategories[category.Raw] = append(rawCategories[category.Raw], machine)
			primaryCategories[category.Primary] = append(primaryCategories[category.Primary], machine)
			secondaryCategories[category.Secondary] = append(secondaryCategories[category.Secondary], machine)
			return nil
		}))

	err = update(db, CategoryRawMachines, rawCategories)
	if err != nil {
		log.Fatal(err)
	}

	err = update(db, CategoryPrimaryMachines, primaryCategories)
	if err != nil {
		log.Fatal(err)
	}

	err = update(db, CategorySecondaryMachines, secondaryCategories)
	if err != nil {
		log.Fatal(err)
	}

}

//GetCategories returns all unqiue categories from boltdb
func GetCategories(db *bolt.DB) []string {
	return GetAllKeys(db, CategoryRawMachines)
}
