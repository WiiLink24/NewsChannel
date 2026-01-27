package tagesschau

import (
	"NewsChannel/news"
	_ "embed"
	"unicode/utf16"
)

type Tagesschau struct {
	oldArticleTitles []string
	news.Source
}

//go:embed logo.jpg
var Logo []byte

func NewTagesschau(oldArticleTitles []string) *Tagesschau {
	return &Tagesschau{
		oldArticleTitles: oldArticleTitles,
	}
}

func (r *Tagesschau) GetLogo() []byte {
	return Logo
}

func (r *Tagesschau) GetCopyright() []uint16 {
	copyrightString := "© ARD-aktuell / tagesschau.de"
	return utf16.Encode([]rune(copyrightString))
}
