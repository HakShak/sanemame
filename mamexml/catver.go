package mamexml

import (
	"log"
	"regexp"
	"strings"
	"time"
)

import "github.com/vaughan0/go-ini"

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
	startTime := time.Now()
	file, err := ini.LoadFile(fileName)
	if err != nil {
		return nil, err
	}

	categories := make(map[string]Category)

	for key, value := range file["Category"] {
		newCategory, err := getCategory(value)
		if err != nil {
			return nil, err
		}

		categories[key] = *newCategory
	}

	log.Printf("Loaded %s in %s", fileName, time.Since(startTime))
	return categories, nil
}
