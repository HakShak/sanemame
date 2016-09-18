package main

import (
	"log"
)

import "github.com/HakShak/sanemame/mamexml"

func check(e error) {
	if e != nil {
		log.Println(e)
	}
}

func main() {
	log.SetFlags(log.Lmicroseconds | log.Lshortfile)
	log.Println("Derp")

	filename := "E:\\temp\\rom\\mame0177.xml"

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
