package db

import (
	"encoding/json"
	"log"

	"github.com/HakShak/sanemame/mamexml"
	"github.com/boltdb/bolt"
)

//NPlayerMachines bucket name
const NPlayerMachines = "nplayer-machines"

//NPlayerTypeMachines bucket name
const NPlayerTypeMachines = "nplayertype-machines"

//NPlayerRawMachines bucket name
const NPlayerRawMachines = "nplayerraw-machines"

func appendMapListOnce(theMap map[string][]string, key string, value string) {
	if !InList(theMap[key], value) {
		theMap[key] = append(theMap[key], value)
	}
}

//UpdateNPlayers populate nplayers data from INI
func UpdateNPlayers(db *bolt.DB, fileName string) {
	nplayers, err := mamexml.LoadNPlayersIni(fileName)
	if err != nil {
		log.Fatal(err)
	}

	nplayerMachines := make(map[string][]string)
	nplayerTypeMachines := make(map[string][]string)
	nplayerRawMachines := make(map[string][]string)

	for machine, nplayer := range nplayers {
		for _, player := range nplayer {
			appendMapListOnce(nplayerMachines, string(player.Players), machine)
			if len(player.PlayerType) > 0 {
				appendMapListOnce(nplayerTypeMachines, player.PlayerType, machine)
			}
			appendMapListOnce(nplayerRawMachines, player.Raw, machine)
		}
	}

	err = db.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte(NPlayerMachines))
		if err != nil {
			return NewCreateBucketError(NPlayerMachines, err)
		}

		for playerCount, machineList := range nplayerMachines {
			machineListBytes, err := json.Marshal(machineList)
			if err != nil {
				return err
			}
			err = bucket.Put([]byte(playerCount), machineListBytes)
			if err != nil {
				return NewPutBucketError(NPlayerMachines, playerCount, string(machineListBytes), err)
			}
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}

	err = db.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte(NPlayerTypeMachines))
		if err != nil {
			return NewCreateBucketError(NPlayerTypeMachines, err)
		}

		for playerType, machineList := range nplayerTypeMachines {
			machineListBytes, err := json.Marshal(machineList)
			if err != nil {
				return err
			}
			err = bucket.Put([]byte(playerType), machineListBytes)
			if err != nil {
				return NewPutBucketError(NPlayerTypeMachines, playerType, string(machineListBytes), err)
			}
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}

	err = db.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte(NPlayerRawMachines))
		if err != nil {
			return NewCreateBucketError(NPlayerRawMachines, err)
		}

		for playerRaw, machineList := range nplayerRawMachines {
			machineListBytes, err := json.Marshal(machineList)
			if err != nil {
				return err
			}
			err = bucket.Put([]byte(playerRaw), machineListBytes)
			if err != nil {
				return NewPutBucketError(NPlayerRawMachines, playerRaw, string(machineListBytes), err)
			}
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
}

//GetNPlayerRawKeys returns all unique nplayer raw keys
func GetNPlayerRawKeys(db *bolt.DB) []string {
	return GetAllKeys(db, NPlayerRawMachines)
}
