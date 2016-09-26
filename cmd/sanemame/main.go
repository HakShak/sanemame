package main

import (
	"log"
)

import "github.com/HakShak/sanemame/mamexml"
import "github.com/HakShak/sanemame/config"

func check(e error) {
	if e != nil {
		log.Println(e)
	}
}

func main() {
	log.SetFlags(0)

	config.SetupConfig()

	categories, err := mamexml.LoadCatverIni("Catver.ini")
	if err != nil {
		log.Fatal(err)
	}

	mature := 0

	for _, value := range categories {
		if value.Mature {
			mature += 1
			log.Printf("%q", value)
		}

	}

	log.Printf("Categories: %d", len(categories))
	log.Printf("Mature: %d", mature)

	nplayers, err := mamexml.LoadNPlayersIni("nplayers.ini")
	if err != nil {
		log.Fatal(err)
	}

	nPlayerCount := 0
	nPlayerUnknown := 0

	for _, value := range nplayers {
		if value[0].Players > 0 {
			nPlayerCount += 1
		} else {
			nPlayerUnknown += 1
		}
	}

	log.Printf("Known NPlayers: %d", nPlayerCount)
	log.Printf("Unknown NPlayers: %d", nPlayerUnknown)

	filename, err := mamexml.GetLatestXmlFile()
	if err != nil {
		log.Fatal(err)
	}

	machines, err := mamexml.Load(filename)
	check(err)

	devices := 0
	bios := 0
	runnable := 0
	mechanical := 0
	clones := 0
	roms := 0
	samples := 0

	for _, m := range machines {
		if m.IsDevice {
			devices += 1
		}
		if m.IsBios {
			bios += 1
		}
		if m.IsRunnable {
			runnable += 1
		}
		if m.IsMechanical {
			mechanical += 1
		}
		if m.CloneOf != "" {
			clones += 1
		}
		if m.RomOf != "" {
			roms += 1
		}
		if m.SampleOf != "" {
			samples += 1
		}
	}

	log.Printf("Machines: %d", len(machines))
	log.Printf("Devices: %d", devices)
	log.Printf("Bios: %d", bios)
	log.Printf("Runnable: %d", runnable)
	log.Printf("Mechanical: %d", mechanical)
	log.Printf("Clones: %d", clones)
	log.Printf("Samples: %d", samples)
	log.Printf("Roms: %d", roms)
	potential := len(machines) - devices - bios - mechanical
	log.Printf("Potential: %d", potential)
	nonClones := runnable - clones
	log.Printf("NonClones: %d", nonClones)
}
