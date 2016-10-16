package db

import (
	"log"

	"github.com/HakShak/sanemame/filetypes/mamexml"
	"github.com/boltdb/bolt"
)

//MameDriverStatusMachines bucket name
const MameDriverStatusMachines = "mamedriverstatus-machines"

//MameBooleanMachines bucket name
const MameBooleanMachines = "mameboolean-machines"

//UpdateMachines populates mamexml into boltdb
func UpdateMachines(db *bolt.DB, fileName string) {
	CreateBuckets(db, []string{
		MameDriverStatusMachines,
		MameBooleanMachines,
	})

	driverMachines := make(map[string][]string)
	propertyMachines := make(map[string][]string)

	err := mamexml.Load(fileName,
		mamexml.EntryLoad(func(machine *mamexml.Machine) error {
			driverMachines[machine.Driver.Status] = append(driverMachines[machine.Driver.Status], machine.Name)

			//It may be an interesting idea to use "reflect" here
			if machine.IsDevice {
				propertyMachines["IsDevice"] = append(propertyMachines["IsDevice"], machine.Name)
			}

			if machine.IsBios {
				propertyMachines["IsBios"] = append(propertyMachines["IsBios"], machine.Name)
			}

			if machine.IsRunnable {
				propertyMachines["IsRunnable"] = append(propertyMachines["IsRunnable"], machine.Name)
			}

			if machine.IsMechanical {
				propertyMachines["IsMechanical"] = append(propertyMachines["IsMechanical"], machine.Name)
			}

			if len(machine.CloneOf) > 0 {
				propertyMachines["IsClone"] = append(propertyMachines["IsClone"], machine.Name)
			}

			if len(machine.SampleOf) > 0 {
				propertyMachines["IsSample"] = append(propertyMachines["IsSample"], machine.Name)
			}

			return nil
		}))
	if err != nil {
		log.Fatal(err)
	}

	err = UpdateStringList(db, MameDriverStatusMachines, driverMachines)
	if err != nil {
		log.Fatal(err)
	}

	err = UpdateStringList(db, MameBooleanMachines, propertyMachines)
	if err != nil {
		log.Fatal(err)
	}
}

//GetDriverStatusMachines returns machines by driver status
func GetDriverStatusMachines(db *bolt.DB) map[string][]string {
	return GetAllLists(db, MameDriverStatusMachines)
}

//GetBooleanMachines returns machines by boolean states
func GetBooleanMachines(db *bolt.DB) map[string][]string {
	return GetAllLists(db, MameBooleanMachines)
}
