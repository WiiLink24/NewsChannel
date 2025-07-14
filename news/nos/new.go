package nos

import (
	_ "embed"
)

//go:embed logo.jpg
var Logo []byte

func NewNos(oldArticleTitles []string) *nos {
	return &nos{
		oldArticleTitles: oldArticleTitles,
	}
}

func (a *nos) GetLogo() []byte {
	return Logo
}
