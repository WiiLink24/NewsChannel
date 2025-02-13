package main

import (
	"NewsChannel/news/reuters"
	"bytes"
	"fmt"
	"github.com/wii-tools/lzx/lz10"
	"hash/crc32"
	"os"
	"testing"
	"time"
)

func makeNews(_t *testing.T, hour int, dayDelta int) {
	var err error
	articles, err = reuters.NewReuters(reuters.UnitedStates).GetArticles()
	if err != nil {
		_t.Fatal(err)
	}

	n := News{}
	n.currentCountryCode = 18
	n.currentLanguageCode = 1

	now := time.Now()
	t := time.Date(now.Year(), now.Month(), now.Day()-dayDelta, hour, 0, 0, 0, time.Local)
	currentTime = int(t.Unix())
	n.currentHour = t.Hour()

	buffer := new(bytes.Buffer)
	// n.GetNewsSource()
	n.MakeHeader()
	n.MakeWiiMenuHeadlines()
	n.ReadNewsCache()
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
	if err != nil {
		_t.Fatal(err)
	}

	err = os.WriteFile(fmt.Sprintf("./v2/1/018/news.dec.%02d", n.currentHour), buffer.Bytes(), 0666)
	err = os.WriteFile(fmt.Sprintf("./v2/1/018/news.bin.%02d", n.currentHour), SignFile(compressed), 0666)
	if err != nil {
		_t.Fatal(err)
	}
}

func TestAllFileGeneration(_t *testing.T) {
	t := time.Now()

	for i := 0; i < t.Hour(); i++ {
		makeNews(_t, i, 0)
	}

	for i := t.Hour(); i < 24; i++ {
		makeNews(_t, i, 1)
	}
}
