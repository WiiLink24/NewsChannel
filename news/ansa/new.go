package ansa

import (
	_ "embed"
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
