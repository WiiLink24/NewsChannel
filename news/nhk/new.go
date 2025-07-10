package nhk

import (
	_ "embed"
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
