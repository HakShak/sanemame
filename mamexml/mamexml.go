package mamexml

import (
	"bufio"
	"encoding/xml"
	"io"
	"log"
	"os"
	"time"
)

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

func Load(filename string) (map[string]Machine, error) {
	startTime := time.Now()
	log.Printf("Loading %s", filename)
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	bufReader := bufio.NewReader(file)

	decoder := xml.NewDecoder(bufReader)

	loaded := make(map[string]Machine)

	for {
		token, err := decoder.Token()
		if err == io.EOF {
			break
		}

		if err != nil {
			return nil, err
		}

		if token == nil {
			break
		}

		switch startElement := token.(type) {
		case xml.StartElement:
			if startElement.Name.Local == "machine" {
				var m Machine
				//Mame xml defaults to everything runnable
				m.IsRunnable = true
				err := decoder.DecodeElement(&m, &startElement)
				if err != nil {
					return nil, err
				}

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

				loaded[m.Name] = m
			}
		}
	}

	log.Printf("%s loaded in %s", filename, time.Since(startTime))
	return loaded, nil
}
