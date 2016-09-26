package mamexml

import (
	"log"
	"regexp"
	"strconv"
	"time"
)

import "github.com/vaughan0/go-ini"

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

func LoadNPlayersIni(fileName string) (map[string][]NPlayer, error) {
	log.Printf("Loading %s", fileName)
	startTime := time.Now()
	file, err := ini.LoadFile(fileName)
	if err != nil {
		return nil, err
	}

	nplayers := make(map[string][]NPlayer)

	for key, value := range file["NPlayers"] {
		newNPlayers, err := getPlayerType(value)
		if err != nil {
			return nil, err
		}

		nplayers[key] = newNPlayers
	}

	log.Printf("Loaded %s in %s", fileName, time.Since(startTime))
	return nplayers, nil
}
