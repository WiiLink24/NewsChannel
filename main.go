package main

import (
	"NewsChannel/news"
	"bytes"
	"encoding/binary"
	"fmt"
	"hash/crc32"
	"io"
	"log"
	"os"
	"time"

	"github.com/wii-tools/lzx/lz10"
)

type News struct {
	Header          Header
	Headlines       []Headlines
	HeadlineText    []uint16
	Topics          []Topic
	Timestamps      []Timestamp
	TopicText       []uint16
	Articles        []Article
	ArticleText     []uint16
	Sources         []Source
	SourcePictures  []byte
	SourceCopyright []byte
	Locations       []Location
	LocationText    []uint16
	Images          []Image
	ImagesData      []byte
	CaptionData     []uint16

	source news.Source

	currentLanguageCode uint8
	currentCountryCode  uint8
	currentHour         int

	// Titles of articles from previous hours. Required for making sure we don't have duplicates.
	oldArticleTitles []string

	// Placeholder for the timestamps for a specific topic.
	timestamps [][]Timestamp

	// Placeholder for locations. Used in order to collect all the used locations without duplicates.
	locations []*news.Location

	articles []news.Article

	// Placeholder for the topics.
	topics []Topic
}

var currentTime = 0

func main() {
	// Load configuration from JSON file instead of having random numbers (added Spain btw)
	config, err := LoadConfig("countries.json")
	checkError(err)

	// Process each country/language combination
	for _, countryConfig := range config.Countries {
		n := News{}
		n.currentCountryCode = countryConfig.CountryCode
		n.currentLanguageCode = countryConfig.LanguageCode

		log.Printf("Processing %s (%s) - Country: %d, Language: %d",
			countryConfig.Name, countryConfig.Language,
			countryConfig.CountryCode, countryConfig.LanguageCode)

		t := time.Now()
		currentTime = int(t.Unix())
		n.currentHour = t.Hour()

		buffer := new(bytes.Buffer)
		n.ReadNewsCache()
		n.setSource(countryConfig.Source)
		n.GetNewsArticles()
		n.MakeHeader()
		n.MakeWiiMenuHeadlines()
		n.MakeArticleTable()
		n.MakeTopicTable()
		n.MakeSourceTable()
		n.WriteNewsCache()
		n.MakeLocationTable()
		n.WriteImages()
		n.Header.Filesize = n.GetCurrentSize()
		n.WriteAll(buffer)

		crcTable := crc32.MakeTable(crc32.IEEE)
		checksum := crc32.Checksum(buffer.Bytes()[12:], crcTable)
		n.Header.CRC32 = checksum

		buffer.Reset()
		n.WriteAll(buffer)

		compressed, err := lz10.Compress(buffer.Bytes())
		checkError(err)

		// If the folder exists we can just continue
		err = os.MkdirAll(fmt.Sprintf("./v2/%d/%03d", n.currentLanguageCode, n.currentCountryCode), os.ModePerm)
		if !os.IsExist(err) {
			checkError(err)
		}

		err = os.WriteFile(fmt.Sprintf("./v2/%d/%03d/news.bin.%02d", n.currentLanguageCode, n.currentCountryCode, n.currentHour), SignFile(compressed), 0666)
		checkError(err)

		log.Printf("Successfully generated news file for %s (%s)", countryConfig.Name, countryConfig.Language)
	}
}

func checkError(err error) {
	if err != nil {
		log.Fatalf("News Channel file generator has encountered a fatal error! Reason: %v\n", err)
	}
}

func Write(writer io.Writer, data any) {
	err := binary.Write(writer, binary.BigEndian, data)
	checkError(err)
}

func (n *News) WriteAll(writer io.Writer) {
	Write(writer, n.Header)
	Write(writer, n.Headlines)
	Write(writer, n.HeadlineText)
	Write(writer, n.Articles)
	Write(writer, n.ArticleText)
	Write(writer, n.Topics)
	Write(writer, n.Timestamps)
	Write(writer, n.TopicText)
	Write(writer, n.Sources)
	Write(writer, n.SourcePictures)
	Write(writer, n.Locations)
	Write(writer, n.LocationText)
	Write(writer, n.Images)
	Write(writer, n.ImagesData)
	Write(writer, n.CaptionData)
}

func (n *News) GetCurrentSize() uint32 {
	buffer := bytes.NewBuffer(nil)
	n.WriteAll(buffer)

	return uint32(buffer.Len())
}
