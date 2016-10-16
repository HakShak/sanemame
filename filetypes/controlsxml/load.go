package controlsxml

import (
	"encoding/xml"

	mxml "github.com/HakShak/sanemame/filetypes/xml"
)

//ControlsName Human description of control
type ControlsName struct {
	Name        string
	Description string
}

//ControlsConstant holds ControlsName
type ControlsConstant struct {
	Name string `xml:"name,attr"`
}

//ControlsControl Why is this structured this way?
type ControlsControl struct {
	Name     string           `xml:"name,attr"`
	Constant ControlsConstant `xml:"constant"`
}

//ControlsPlayer I give up
type ControlsPlayer struct {
	Controls []ControlsControl `xml:"controls>control"`
}

//ControlsGame all control related to a game element
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

//GetControlsNames find all names for a certain constant
func GetControlsNames(controlGame ControlsGame) []ControlsName {
	var result []ControlsName
	for _, playerControl := range controlGame.PlayerControls {
		for _, control := range playerControl.Controls {
			result = append(result, ControlsName{control.Constant.Name, control.Name})
		}
	}
	return result
}

//EntryRead callback signature
type EntryRead func(controls *ControlsGame) error

//Load process all control elements
func Load(fileName string, callback EntryRead) error {
	err := mxml.Load(fileName, mxml.ElementLoad(
		func(decoder *xml.Decoder, element *xml.StartElement) error {
			if element.Name.Local == "game" {
				var cg ControlsGame
				err := decoder.DecodeElement(&cg, element)
				if err != nil {
					return err
				}

				err = callback(&cg)
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
