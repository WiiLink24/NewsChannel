package rtve

import (
	"NewsChannel/news"
	_ "embed"
)

type RTVE struct {
	oldArticleTitles []string
	news.Source
}

//go:embed logo.jpg
var Logo []byte

func NewRTVE(oldArticleTitles []string) *RTVE {
	return &RTVE{
		oldArticleTitles: oldArticleTitles,
	}
}

func (r *RTVE) GetLogo() []byte {
	return Logo
}

func (r *RTVE) GetCopyright() string {
	return "© RTVE. All rights reserved."
}
