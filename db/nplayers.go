package db

import (
	"log"

	"github.com/HakShak/sanemame/filetypes/nplayersini"
	"github.com/boltdb/bolt"
)

//NPlayerMachines bucket name
const NPlayerMachines = "nplayer-machines"

//NPlayerTypeMachines bucket name
const NPlayerTypeMachines = "nplayertype-machines"

//NPlayerRawMachines bucket name
const NPlayerRawMachines = "nplayerraw-machines"

//UpdateNPlayers populate nplayers data from INI
func UpdateNPlayers(db *bolt.DB, fileName string) {
	CreateBuckets(db, []string{
		NPlayerMachines,
		NPlayerRawMachines,
		NPlayerTypeMachines,
	})

	nplayerMachines := make(map[string][]string)
	nplayerTypeMachines := make(map[string][]string)
	nplayerRawMachines := make(map[string][]string)

	err := nplayersini.Load(fileName, nplayersini.EntryRead(
		func(machine string, nplayer *nplayersini.NPlayer) error {
			nplayerRawMachines[nplayer.Raw] = append(nplayerRawMachines[nplayer.Raw], machine)
			nplayerTypeMachines[nplayer.PlayerType] = append(nplayerTypeMachines[nplayer.PlayerType], machine)
			nplayerMachines[string(nplayer.Players)] = append(nplayerMachines[string(nplayer.Players)], machine)
			return nil
		}))
	if err != nil {
		log.Fatal(err)
	}

	err = UpdateStringList(db, NPlayerMachines, nplayerMachines)
	if err != nil {
		log.Fatal(err)
	}

	err = UpdateStringList(db, NPlayerRawMachines, nplayerRawMachines)
	if err != nil {
		log.Fatal(err)
	}

	err = UpdateStringList(db, NPlayerTypeMachines, nplayerTypeMachines)
	if err != nil {
		log.Fatal(err)
	}
}

//GetNPlayerRawKeys returns all unique nplayer raw keys
func GetNPlayerRawKeys(db *bolt.DB) []string {
	return GetAllKeys(db, NPlayerRawMachines)
}

//GetNPlayerMachines returns machines by player count
func GetNPlayerMachines(db *bolt.DB) map[string][]string {
	return GetAllLists(db, NPlayerMachines)
}

//GetNPlayerRawMachines returns machines by raw nplayer notation
func GetNPlayerRawMachines(db *bolt.DB) map[string][]string {
	return GetAllLists(db, NPlayerRawMachines)
}

//GetNPlayerTypeMachines returns machines by player type
func GetNPlayerTypeMachines(db *bolt.DB) map[string][]string {
	return GetAllLists(db, NPlayerTypeMachines)
}
