package nos

import (
	_ "embed"
	"fmt"
	"strconv"
	"time"
	"unicode/utf16"
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

func (a *nos) GetCopyright() []uint16 {
	copyrightString := fmt.Sprintf("© NOS %s", strconv.Itoa(time.Now().Year()))
	return utf16.Encode([]rune(copyrightString))
}
