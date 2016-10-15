package ini

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"

	pb "gopkg.in/cheggaaa/pb.v1"
)

var (
	sectionRegex = regexp.MustCompile(`^\[(.*)\]$`)
	assignRegex  = regexp.MustCompile(`^([^=]+)=(.*)$`)
)

//EntryRead callback to process each load
type EntryRead func(section, key, value string) error

//Load ini with callback for processing
func Load(fileName string, callback EntryRead) error {
	log.Printf("Loading %s", fileName)
	file, err := os.Open(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	info, err := file.Stat()
	if err != nil {
		return err
	}

	bar := pb.New(int(info.Size())).SetUnits(pb.U_BYTES).SetWidth(80).Start()
	bar.ShowSpeed = true

	bufReader := bufio.NewReader(file)

	proxyReader := bar.NewProxyReader(bufReader)

	lineNumber := 0
	scanner := bufio.NewScanner(proxyReader)
	currentSection := ""

	for scanner.Scan() {
		lineNumber++
		section, key, value, err := parseLine(scanner.Text(), lineNumber)
		if err != nil {
			return err
		}

		if len(section) > 0 {
			currentSection = section
			continue
		}

		//Only trigger callbacks on something useful
		if len(key) == 0 || len(value) == 0 {
			continue
		}

		err = callback(currentSection, key, value)
		if err != nil {
			return err
		}
	}

	bar.Finish()

	return nil
}

func parseLine(line string, lineNumber int) (section string, key string, value string, err error) {
	line = strings.TrimSpace(line)
	if len(line) == 0 {
		// Skip blank lines
		return
	}
	if line[0] == ';' || line[0] == '#' {
		// Skip comments
		return
	}

	if groups := assignRegex.FindStringSubmatch(line); groups != nil {
		key, value := groups[1], groups[2]
		key, value = strings.TrimSpace(key), strings.TrimSpace(value)
		return "", key, value, nil
	} else if groups := sectionRegex.FindStringSubmatch(line); groups != nil {
		name := strings.TrimSpace(groups[1])
		section = name
		return
	} else {
		msg := fmt.Sprintf("Syntax error on line %d: %s", lineNumber, line)
		return "", "", "", errors.New(msg)
	}
}
