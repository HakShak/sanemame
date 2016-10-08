package db

import (
	"encoding/json"
	"github.com/HakShak/sanemame/mamexml"
	"github.com/boltdb/bolt"
	"log"
)

const CategoryMachines = "category-machines"

func UpdateCategories(db *bolt.DB, fileName string) {
	categories, err := mamexml.LoadCatverIni(fileName)
	if err != nil {
		log.Fatal(err)
	}

	categorySet := make(map[string][]string)

	for machine, category := range categories {
		categoryKey := category.Primary
		if len(category.Secondary) > 0 {
			categoryKey += "-" + category.Secondary
		}
		categorySet[categoryKey] = append(categorySet[categoryKey], machine)
	}

	err = db.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte(CategoryMachines))
		if err != nil {
			return err
		}

		for categoryKey, machineList := range categorySet {
			machineListBytes, err := json.Marshal(machineList)
			if err != nil {
				return err
			}
			err = bucket.Put([]byte(categoryKey), machineListBytes)
			if err != nil {
				return err
			}
		}
		return nil
	})

	if err != nil {
		log.Fatal(err)
	}
}

func GetCategories(db *bolt.DB) []string {
	return GetAllKeys(db, CategoryMachines)
}
