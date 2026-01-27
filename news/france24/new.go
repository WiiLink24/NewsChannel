package france24

import (
	_ "embed"
	"fmt"
	"strconv"
	"time"
	"unicode/utf16"
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

func (a *france24) GetCopyright() []uint16 {
	copyrightString := fmt.Sprintf("© %s Copyright France 24 - Tous droits réservés.", strconv.Itoa(time.Now().Year()))
	return utf16.Encode([]rune(copyrightString))
}
