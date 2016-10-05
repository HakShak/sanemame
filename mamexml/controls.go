package mamexml

import (
	"bufio"
	"encoding/xml"
	"io"
	"log"
	"os"
)

import "gopkg.in/cheggaaa/pb.v1"

type ControlsConstant struct {
	Name string `xml:"name,attr"`
}

type ControlsControl struct {
	Name     string           `xml:"name,attr"`
	Constant ControlsConstant `xml:"constant"`
}

type ControlsPlayer struct {
	Controls []ControlsControl `xml:"controls>control"`
}

type ControlsGame struct {
	RomName        string           `xml:"romname,attr"`
	GameName       string           `xml:"gamename,attr"`
	Players        int              `xml:"numPlayers,attr"`
	Alternating    bool             `xml:"alternating,attr"`
	Mirrored       bool             `xml:"mirrored,attr"`
	UsesService    bool             `xml:"usesService,attr"`
	Tilt           bool             `xml:"tilt,attr"`
	Cocktail       bool             `xml:"cocktail,attr"`
	PlayerControls []ControlsPlayer `xml:"player"`
}

func LoadControlsXml(filename string) (map[string]ControlsGame, error) {
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

	bar.Finish()
	return loaded, nil
}
