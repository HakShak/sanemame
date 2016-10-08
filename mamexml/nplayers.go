package mamexml

import (
	"bufio"
	"log"
	"os"
	"regexp"
	"strconv"
)

import "github.com/vaughan0/go-ini"
import "gopkg.in/cheggaaa/pb.v1"

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

	file, err := os.Open(fileName)
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

	iniFile, err := ini.Load(proxyReader)
	if err != nil {
		return nil, err
	}

	nplayers := make(map[string][]NPlayer)

	for key, value := range iniFile["NPlayers"] {
		newNPlayers, err := getPlayerType(value)
		if err != nil {
			return nil, err
		}

		nplayers[key] = newNPlayers
	}

	bar.Finish()
	return nplayers, nil
}
