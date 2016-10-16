package catverini

import (
	"regexp"
	"strings"

	"github.com/HakShak/sanemame/filetypes/ini"
)

//Category parsed struct from catver.ini
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

//EntryRead callback to process each load
type EntryRead func(machine string, category *Category) error

//Load ini with callback for processing
func Load(fileName string, callback EntryRead) error {
	err := ini.Load(fileName,
		ini.EntryRead(func(section, key, value string) error {
			if section != "Category" {
				return nil
			}

			category, err := getCategory(value)
			if err != nil {
				return err
			}

			err = callback(key, category)
			if err != nil {
				return err
			}
			return nil
		}))
	if err != nil {
		return err
	}
	return nil
}
