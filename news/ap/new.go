package ap

import (
	_ "embed"
	"fmt"
	"strconv"
	"time"
	"unicode/utf16"
)

type AP struct {
	oldArticleTitles []string
}

//go:embed logo.jpg
var Logo []byte

func NewAP(oldArticleTitles []string) *AP {
	return &AP{
		oldArticleTitles: oldArticleTitles,
	}
}

func (a *AP) GetLogo() []byte {
	return Logo
}

func (a *AP) GetCopyright() []uint16 {
	copyrightString := fmt.Sprintf("Copyright %s The Associated Press. All rights reserved.", strconv.Itoa(time.Now().Year()))
	return utf16.Encode([]rune(copyrightString))
}
