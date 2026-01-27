package nhk

import (
	_ "embed"
	"unicode/utf16"
)

//go:embed logo.jpg
var Logo []byte

func NewNHK(oldArticleTitles []string) *nhk {
	return &nhk{
		oldArticleTitles: oldArticleTitles,
	}
}

func (a *nhk) GetLogo() []byte {
	return Logo
}

func (a *nhk) GetCopyright() []uint16 {
	copyrightString := "Copyright NHK (Japan Broadcasting Corporation). All rights reserved."
	return utf16.Encode([]rune(copyrightString))
}
