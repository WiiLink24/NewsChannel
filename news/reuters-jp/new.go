package reutersjp

import (
	"NewsChannel/news"
	_ "embed"
	"fmt"
	"strconv"
	"time"
	"unicode/utf16"
)

type ReutersJP struct {
	oldArticleTitles []string
	news.Source
}

//go:embed logo.jpg
var Logo []byte

func NewReuters(oldArticleTitles []string) *ReutersJP {
	return &ReutersJP{
		oldArticleTitles: oldArticleTitles,
	}
}

func (r *ReutersJP) GetLogo() []byte {
	return Logo
}

func (r *ReutersJP) GetCopyright() []uint16 {
	copyrightString := fmt.Sprintf("© %s Reuters. All rights reserved", strconv.Itoa(time.Now().Year()))
	return utf16.Encode([]rune(copyrightString))
}
