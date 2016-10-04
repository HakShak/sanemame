package mamexml

import (
	"bufio"
	"log"
	"os"
	"regexp"
	"strings"
)

import "github.com/vaughan0/go-ini"
import "gopkg.in/cheggaaa/pb.v1"

type Category struct {
	Raw       string
	Primary   string
	Secondary string
	Mature    bool
}

func getCategory(value string) (*Category, error) {
	re := regexp.MustCompile("[^\\/\\*]+")
	matches := re.FindAllString(value, -1)
	if len(matches) == 0 {
		return nil, nil
	}

	mature := false
	primary := ""
	secondary := ""

	for index, match := range matches {
		match = strings.TrimSpace(match)
		if match == "Mature" {
			mature = true
			//This seems to always be the last "tag"
			break
		}

		if index == 0 {
			primary = match
			continue
		}

		if index == 1 {
			secondary = match
			continue
		}
	}

	result := Category{value, primary, secondary, mature}

	return &result, nil
}

func LoadCatverIni(fileName string) (map[string]Category, error) {
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

	categories := make(map[string]Category)

	for key, value := range iniFile["Category"] {
		newCategory, err := getCategory(value)
		if err != nil {
			return nil, err
		}

		categories[key] = *newCategory
	}

	bar.Finish()

	return categories, nil
}
