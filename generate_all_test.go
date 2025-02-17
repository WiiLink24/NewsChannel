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
	countries := []uint8{18, 49, 110}
	for _, country := range countries {
		n := News{}
		n.currentCountry = reuters.GetCountry(country)
		n.currentCountryCode = country
		n.currentLanguageCode = 1

		now := time.Now()
		t := time.Date(now.Year(), now.Month(), now.Day()-dayDelta, hour, 0, 0, 0, time.Local)
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
		if err != nil {
			_t.Fatal(err)
		}

		// If the folder exists we can just continue
		err = os.MkdirAll(fmt.Sprintf("./v2/%d/%03d", n.currentLanguageCode, n.currentCountryCode), os.ModePerm)
		if !os.IsExist(err) {
			if err != nil {
				_t.Fatal(err)
			}
		}

		err = os.WriteFile(fmt.Sprintf("./v2/%d/%03d/news.bin.%02d", n.currentLanguageCode, n.currentCountryCode, n.currentHour), SignFile(compressed), 0666)
		if err != nil {
			_t.Fatal(err)
		}
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
