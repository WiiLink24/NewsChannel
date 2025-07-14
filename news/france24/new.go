package france24

import (
	_ "embed"
)

//go:embed logo.jpg
var Logo []byte

func NewFrance24(oldArticleTitles []string) *france24 {
	return &france24{
		oldArticleTitles: oldArticleTitles,
	}
}

func (a *france24) GetLogo() []byte {
	return Logo
}
