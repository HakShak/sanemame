package db

import (
	"encoding/json"
	"github.com/HakShak/sanemame/mamexml"
	"github.com/boltdb/bolt"
	"log"
)

const ControlMachines = "control-machines"
const ControlNames = "control-names"

func checkList(list []string, key string) bool {
	for _, check := range list {
		if check == key {
			return true
		}
	}
	return false
}

func UpdateControls(db *bolt.DB, fileName string) {
	controls, err := mamexml.LoadControlsXml(fileName)
	if err != nil {
		log.Fatal(err)
	}

	controlNames := make(map[string][]string)
	controlMachines := make(map[string][]string)

	for machine, controls := range controls {
		for _, control := range mamexml.GetControlsNames(controls) {
			if !checkList(controlNames[control.Name], control.Description) {
				controlNames[control.Name] = append(controlNames[control.Name], control.Description)
			}

			controlMachines[control.Name] = append(controlMachines[control.Name], machine)
		}
	}

	err = db.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte(ControlNames))
		if err != nil {
			return err
		}

		for controlKey, nameList := range controlNames {
			nameListBytes, err := json.Marshal(nameList)
			if err != nil {
				return err
			}
			err = bucket.Put([]byte(controlKey), nameListBytes)
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}

	err = db.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte(ControlMachines))
		if err != nil {
			return err
		}

		for controlKey, machineList := range controlMachines {
			machineListBytes, err := json.Marshal(machineList)
			if err != nil {
				return err
			}
			err = bucket.Put([]byte(controlKey), machineListBytes)
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

func GetControls(db *bolt.DB) []string {
	return GetAllKeys(db, ControlNames)
}

func GetControlNames(db *bolt.DB) map[string]string {
	return GetAll(db, ControlNames)
}
