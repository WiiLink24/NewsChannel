package rtve

import (
	"NewsChannel/news"
	_ "embed"
	"fmt"
	"strconv"
	"time"
	"unicode/utf16"
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

func (r *RTVE) GetCopyright() []uint16 {
	copyrightString := fmt.Sprintf(" © Corporación de Radio y Televisión Española %s", strconv.Itoa(time.Now().Year()))
	return utf16.Encode([]rune(copyrightString))
}
