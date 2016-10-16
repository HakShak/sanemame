package xml

import (
	"bufio"
	"encoding/xml"
	"io"
	"log"
	"os"

	pb "gopkg.in/cheggaaa/pb.v1"
)

//ElementLoad callback to process elements
type ElementLoad func(decoder *xml.Decoder, element *xml.StartElement) error

//Load load file with processing callback
func Load(filename string, callback ElementLoad) error {
	log.Printf("Loading %s", filename)
	file, err := os.Open(filename)
	if err != nil {
		return err
	}

	info, err := file.Stat()
	if err != nil {
		return err
	}

	bar := pb.New(int(info.Size())).SetUnits(pb.U_BYTES).SetWidth(80).Start()
	bar.ShowSpeed = true

	bufReader := bufio.NewReader(file)

	proxyReader := bar.NewProxyReader(bufReader)

	decoder := xml.NewDecoder(proxyReader)

	for {
		token, err := decoder.Token()
		if err == io.EOF {
			break
		}

		if err != nil {
			return err
		}

		if token == nil {
			break
		}

		switch startElement := token.(type) {
		case xml.StartElement:
			err := callback(decoder, &startElement)
			if err != nil {
				return err
			}
		}
	}

	bar.Finish()
	return nil
}
