package nplayersini

import (
	"regexp"
	"strconv"

	"github.com/HakShak/sanemame/filetypes/ini"
)

//NPlayer Represents player counts and types
type NPlayer struct {
	Raw        string
	Players    int
	PlayerType string
}

func getPlayerType(value string) ([]NPlayer, error) {
	re := regexp.MustCompile("(\\d)P (\\w+)")
	matches := re.FindAllStringSubmatch(value, -1)
	if len(matches) == 0 {
		return []NPlayer{{value, 0, ""}}, nil
	}

	results := []NPlayer{}

	for _, match := range matches {
		players, err := strconv.Atoi(match[1])
		if err != nil {
			return nil, err
		}

		newNPlayer := NPlayer{match[0], players, match[2]}
		results = append(results, newNPlayer)
	}

	return results, nil
}

//EntryRead callback to process each load
type EntryRead func(machine string, nplayer *NPlayer) error

//Load ini with callback for processing
func Load(fileName string, callback EntryRead) error {
	ini.Load(fileName,
		ini.EntryRead(func(section, key, value string) error {
			if section != "NPlayers" {
				return nil
			}

			newNPlayers, err := getPlayerType(value)
			if err != nil {
				return err
			}

			for _, nPlayer := range newNPlayers {
				err = callback(key, &nPlayer)
				if err != nil {
					return err
				}
			}

			return nil
		}))

	return nil
}
