package reuters

import (
	"NewsChannel/news"
	_ "embed"
)

// Country represents a country in the Reuters API
type Country string

const (
	UnitedStates  = "us"
	Canada        = "canada"
	Germany       = "germany"
	France        = "france"
	UnitedKingdom = "uk"
	Australia     = "australia"
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
	case 65:
		return Australia
	case 78:
		return Germany
	case 77:
		return France
	case 18:
		return Canada
	case 110:
		return UnitedKingdom
	default:
		return UnitedStates
	}
}

func (r *Reuters) GetLogo() []byte {
	return Logo
}
