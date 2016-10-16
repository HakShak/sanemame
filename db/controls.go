package db

import (
	"log"

	"github.com/HakShak/sanemame/filetypes/controlsxml"
	"github.com/boltdb/bolt"
)

//ControlMachines bucket name
const ControlMachines = "control-machines"

//ControlNames bucket name
const ControlNames = "control-names"

//UpdateControls Populates controls from XML into boltdb
func UpdateControls(db *bolt.DB, fileName string) {
	CreateBuckets(db, []string{
		ControlMachines,
		ControlNames,
	})

	controlNames := make(map[string][]string)
	controlMachines := make(map[string][]string)

	err := controlsxml.Load(fileName,
		controlsxml.EntryRead(func(controls *controlsxml.ControlsGame) error {
			for _, control := range controlsxml.GetControlsNames(*controls) {
				if !InList(controlNames[control.Name], control.Description) {
					controlNames[control.Name] = append(controlNames[control.Name], control.Description)
				}
				controlMachines[control.Name] = append(controlMachines[control.Name], controls.RomName)
			}
			return nil
		}))
	if err != nil {
		log.Fatal(err)
	}

	err = UpdateStringList(db, ControlNames, controlNames)
	if err != nil {
		log.Fatal(err)
	}

	err = UpdateStringList(db, ControlMachines, controlMachines)
	if err != nil {
		log.Fatal(err)
	}
}

//GetControls returns unique control constants
func GetControls(db *bolt.DB) []string {
	return GetAllKeys(db, ControlNames)
}

//GetControlNames returns unique control constants with freesytle names
func GetControlNames(db *bolt.DB) map[string]string {
	return GetAll(db, ControlNames)
}

//GetControlMachines returns all machines by control
func GetControlMachines(db *bolt.DB) map[string][]string {
	return GetAllLists(db, ControlMachines)
}
