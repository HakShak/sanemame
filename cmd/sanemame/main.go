package main

import "log"

import "bufio"
import "os"

import "encoding/xml"

func check(e error) {
	if e != nil {
		log.Println(e)
	}
}

type Driver struct {
	Status string `xml:"status,attr"`
}

type Input struct {
	Players int `xml:"players,attr"`
}

type Machine struct {
	Name         string `xml:"name,attr"`
	Description  string `xml:"description"`
	Driver       Driver `xml:"driver"`
	Input        Input  `xml:"input"`
	IsDevice     bool
	IsBios       bool
	IsRunnable   bool
	IsMechanical bool
	CloneOf      string `xml:"cloneof,attr"`
	RomOf        string `xml:"romof,attr"`
	SampleOf     string `xml:"sampleof,attr"`
}

func main() {
	log.SetFlags(log.Lmicroseconds | log.Lshortfile)
	log.Println("Derp")
	file, err := os.Open("E:\\temp\\rom\\mame0177.xml")
	check(err)

	bufReader := bufio.NewReader(file)

	decoder := xml.NewDecoder(bufReader)

	machines := 0
	devices := 0
	bios := 0
	runnable := 0
	mechanical := 0
	clones := 0
	samples := 0
	roms := 0

	for {
		token, err := decoder.Token()
		check(err)
		if token == nil {
			break
		}

		switch startElement := token.(type) {
		case xml.StartElement:
			if startElement.Name.Local == "machine" {
				var m Machine
				//Mame xml defaults to everything runnable
				m.IsRunnable = true
				decoder.DecodeElement(&m, &startElement)

				for _, attribute := range startElement.Attr {
					switch attribute.Name.Local {
					case "isdevice":
						if attribute.Value == "yes" {
							m.IsDevice = true
						}
					case "isbios":
						if attribute.Value == "yes" {
							m.IsBios = true
						}
					case "ismechanical":
						if attribute.Value == "yes" {
							m.IsMechanical = true
						}
					case "runnable":
						if attribute.Value == "no" {
							m.IsRunnable = false
						}
					}
				}

				machines += 1
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
		}
	}

	log.Printf("Machines: %d", machines)
	log.Printf("Devices: %d", devices)
	log.Printf("Bios: %d", bios)
	log.Printf("Runnable: %d", runnable)
	log.Printf("Mechanical: %d", mechanical)
	log.Printf("Clones: %d", clones)
	log.Printf("Samples: %d", samples)
	log.Printf("Roms: %d", roms)
	potential := machines - devices - bios - mechanical
	log.Printf("Potential: %d", potential)
	nonClones := runnable - clones
	log.Printf("NonClones: %d", nonClones)
}
