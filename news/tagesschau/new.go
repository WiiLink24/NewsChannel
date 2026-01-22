package tagesschau

import (
	"NewsChannel/news"
	_ "embed"
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
