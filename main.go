package main

import (
	"NewsChannel/news"
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/wii-tools/lzx/lz10"
	"hash/crc32"
	"io"
	"log"
	"os"
	"time"
)

type News struct {
	Header         Header
	Headlines      []Headlines
	HeadlineText   []uint16
	Topics         []Topic
	Timestamps     []Timestamp
	TopicText      []uint16
	Articles       []Article
	ArticleText    []uint16
	Sources        []Source
	SourcePictures []byte
	Locations      []Location
	LocationText   []uint16
	Images         []Image
	ImagesData     []byte
	CaptionData    []uint16

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
	// TODO: All Countries!!!!!
	n := News{}
	n.currentCountryCode = 18
	n.currentLanguageCode = 1

	t := time.Now()
	currentTime = int(t.Unix())
	n.currentHour = t.Hour()

	buffer := new(bytes.Buffer)
	n.ReadNewsCache()
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

	err = os.WriteFile(fmt.Sprintf("./v2/1/018/news.dec.%02d", n.currentHour), buffer.Bytes(), 0666)
	err = os.WriteFile(fmt.Sprintf("./v2/1/018/news.bin.%02d", n.currentHour), SignFile(compressed), 0666)
	checkError(err)
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
