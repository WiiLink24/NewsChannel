package reuters

import (
	"NewsChannel/news"
	_ "embed"
)

// Country represents a country in the Reuters API
type Country string

const (
	Japan         = "japan"
	Brazil        = "brazil"
	Canada        = "canada"
	Mexico        = "mexico"
	Peru          = "peru"
	UnitedStates  = "us"
	Uruguay       = "uruguay"
	Australia     = "australia"
	France        = "france"
	Germany       = "germany"
	Mozambique    = "mozambique"
	Poland        = "poland"
	Russia        = "russia"
	SouthAfrica   = "south-africa"
	UnitedKingdom = "uk"
	Azerbaijan    = "azerbaijan"
	Sudan         = "sudan"
	Taiwan        = "taiwan"
	SouthKorea    = "south-korea"
	Singapore     = "singapore"
	China         = "china"
	India         = "india"
	Syria         = "syria"
)

type Reuters struct {
	country          Country
	oldArticleTitles []string
	news.Source
}

//go:embed logo.jpg
var Logo []byte

func NewReuters(oldArticleTitles []string, countryCode uint8) *Reuters {
	return &Reuters{
		oldArticleTitles: oldArticleTitles,
		country:          getCountry(countryCode),
	}
}

func getCountry(code uint8) Country {
	switch code {
	case 1:
		return Japan
	case 16:
		return Brazil
	case 18:
		return Canada
	case 36:
		return Mexico
	case 42:
		return Peru
	case 50:
		return Uruguay
	case 65:
		return Australia
	case 77:
		return France
	case 78:
		return Germany
	case 92:
		return Mozambique
	case 97:
		return Poland
	case 100:
		return Russia
	case 104:
		return SouthAfrica
	case 110:
		return UnitedKingdom
	case 113:
		return Azerbaijan
	case 118:
		return Sudan
	case 128:
		return Taiwan
	case 136:
		return SouthKorea
	case 153:
		return Singapore
	case 160:
		return China
	case 169:
		return India
	case 175:
		return Syria
	default:
		return UnitedStates
	}
}

func (r *Reuters) GetLogo() []byte {
	return Logo
}
