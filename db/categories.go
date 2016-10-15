package db

import (
	"fmt"
	"log"

	"github.com/HakShak/sanemame/filetypes/catverini"
	"github.com/boltdb/bolt"
)

//CategoryRawMachines bucket name
const CategoryRawMachines = "categoryraw-machines"

//CategoryPrimaryMachines bucket name
const CategoryPrimaryMachines = "categoryprimary-machines"

//CategorySecondaryMachines bucket name
const CategorySecondaryMachines = "categorysecondary-machines"

//MatureCategory key name
const MatureCategory = "Mature"

//UpdateCategories populates categories from INI
func UpdateCategories(db *bolt.DB, fileName string) {
	CreateBuckets(db, []string{
		CategoryPrimaryMachines,
		CategoryRawMachines,
		CategorySecondaryMachines,
	})

	rawCategories := make(map[string][]string)
	primaryCategories := make(map[string][]string)
	secondaryCategories := make(map[string][]string)

	err := catverini.Load(fileName, catverini.EntryRead(
		func(machine string, category *catverini.Category) error {
			raw := category.Primary
			if len(category.Secondary) > 0 {
				raw = fmt.Sprintf("%s / %s", category.Primary, category.Secondary)
			}
			rawCategories[raw] = append(rawCategories[raw], machine)
			primaryCategories[category.Primary] = append(primaryCategories[category.Primary], machine)
			secondaryCategories[category.Secondary] = append(secondaryCategories[category.Secondary], machine)

			if category.Mature {
				rawCategories[MatureCategory] = append(rawCategories[MatureCategory], machine)
				primaryCategories[MatureCategory] = append(primaryCategories[MatureCategory], machine)
			}
			return nil
		}))
	if err != nil {
		log.Fatal(err)
	}

	err = UpdateStringList(db, CategoryRawMachines, rawCategories)
	if err != nil {
		log.Fatal(err)
	}

	err = UpdateStringList(db, CategoryPrimaryMachines, primaryCategories)
	if err != nil {
		log.Fatal(err)
	}

	err = UpdateStringList(db, CategorySecondaryMachines, secondaryCategories)
	if err != nil {
		log.Fatal(err)
	}

}

//GetRawCategories returns all unqiue categories from boltdb
func GetRawCategories(db *bolt.DB) []string {
	return GetAllKeys(db, CategoryRawMachines)
}

//GetRawCategoryMachines returns all machines data for raw categories
func GetRawCategoryMachines(db *bolt.DB) map[string][]string {
	return GetAllLists(db, CategoryRawMachines)
}

//GetPrimaryCategories return all unique primary categories
func GetPrimaryCategories(db *bolt.DB) []string {
	return GetAllKeys(db, CategoryPrimaryMachines)
}

//GetPrimaryCategoryMachines returns all machines data for primary categories
func GetPrimaryCategoryMachines(db *bolt.DB) map[string][]string {
	return GetAllLists(db, CategoryPrimaryMachines)
}

//GetSecondaryCategories return all unique primary categories
func GetSecondaryCategories(db *bolt.DB) []string {
	return GetAllKeys(db, CategorySecondaryMachines)
}

//GetSecondaryCategoryMachines returns all machines data for secondary categories
func GetSecondaryCategoryMachines(db *bolt.DB) map[string][]string {
	return GetAllLists(db, CategorySecondaryMachines)
}
