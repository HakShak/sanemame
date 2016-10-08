package mamexml

import (
	"bufio"
	"encoding/xml"
	"io"
	"log"
	"os"
)

import "gopkg.in/cheggaaa/pb.v1"

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

func GetMachines() (interface{}, error) {

	return nil, nil
}

func Load(filename string) (map[string]Machine, error) {
	log.Printf("Loading %s", filename)
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	info, err := file.Stat()
	if err != nil {
		return nil, err
	}

	bar := pb.New(int(info.Size())).SetUnits(pb.U_BYTES).SetWidth(80).Start()
	bar.ShowSpeed = true

	bufReader := bufio.NewReader(file)

	proxyReader := bar.NewProxyReader(bufReader)

	decoder := xml.NewDecoder(proxyReader)

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

	bar.Finish()
	return loaded, nil
}
