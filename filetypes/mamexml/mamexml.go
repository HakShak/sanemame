package mamexml

import (
	"encoding/xml"

	mxml "github.com/HakShak/sanemame/filetypes/xml"
)

//Driver status of emulation
type Driver struct {
	Status string `xml:"status,attr"`
}

//Input number of players
type Input struct {
	Players int `xml:"players,attr"`
}

//Machine basic element
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

//EntryLoad callback
type EntryLoad func(machine *Machine) error

//Load load all data through callback
func Load(fileName string, callback EntryLoad) error {
	err := mxml.Load(fileName, mxml.ElementLoad(func(decoder *xml.Decoder, element *xml.StartElement) error {
		if element.Name.Local == "machine" {
			var m Machine
			//Mame xml defaults to everything runnable
			m.IsRunnable = true
			err := decoder.DecodeElement(&m, element)
			if err != nil {
				return err
			}

			for _, attribute := range element.Attr {
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

			err = callback(&m)
			if err != nil {
				return err
			}
		}

		return nil
	}))
	if err != nil {
		return err
	}

	return nil
}
