package welt

import (
	_ "embed"
)

//go:embed logo.jpg
var Logo []byte

func NewWelt(oldArticleTitles []string) *welt {
	return &welt{
		oldArticleTitles: oldArticleTitles,
	}
}

func (a *welt) GetLogo() []byte {
	return Logo
}

func (a *welt) GetCopyright() string {
	return "© WELT. All rights reserved."
}
