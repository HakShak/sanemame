package mamexml

import (
	"bufio"
	"encoding/xml"
	"io"
	"log"
	"os"
	"time"
)

type ControlsGame struct {
	RomName     string `xml:"romname,attr"`
	GameName    string `xml:"gamename,attr"`
	Players     int    `xml:"numPlayers,attr"`
	Alternating bool   `xml:"alternating,attr"`
	Mirrored    bool   `xml:"mirrored,attr"`
	UsesService bool   `xml:"usesService,attr"`
	Tilt        bool   `xml:"tilt,attr"`
	Cocktail    bool   `xml:"cocktail,attr"`
}

func LoadControlsXml(filename string) (map[string]ControlsGame, error) {
	startTime := time.Now()
	log.Printf("Loading %s", filename)
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	bufReader := bufio.NewReader(file)

	decoder := xml.NewDecoder(bufReader)

	loaded := make(map[string]ControlsGame)

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
			if startElement.Name.Local == "game" {
				var cg ControlsGame
				err := decoder.DecodeElement(&cg, &startElement)
				if err != nil {
					return nil, err
				}

				loaded[cg.RomName] = cg
			}
		}
	}

	log.Printf("%s loaded in %s", filename, time.Since(startTime))
	return loaded, nil
}
