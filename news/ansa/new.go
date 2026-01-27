package ansa

import (
	_ "embed"
	"fmt"
	"strconv"
	"time"
	"unicode/utf16"
)

type ANSA struct {
	oldArticleTitles []string
}

//go:embed logo.jpg
var Logo []byte

func NewAnsa(oldArticleTitles []string) *ANSA {
	return &ANSA{
		oldArticleTitles: oldArticleTitles,
	}
}

func (a *ANSA) GetLogo() []byte {
	return Logo
}

func (a *ANSA) GetCopyright() []uint16 {
	copyrightString := fmt.Sprintf("Copyright %s © ANSA\nTutti i diritti riservati", strconv.Itoa(time.Now().Year()))
	return utf16.Encode([]rune(copyrightString))
}
