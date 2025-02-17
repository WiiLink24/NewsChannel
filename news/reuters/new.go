package reuters

import "NewsChannel/news"

// Country represents a country in the Reuters API
type Country string

const (
	UnitedStates  = "us"
	Canada        = "canada"
	Germany       = "germany"
	France        = "france"
	UnitedKingdom = "uk"
)

type Reuters struct {
	country          Country
	oldArticleTitles []string
	news.Source
}

func NewReuters(oldArticleTitles []string, country Country) *Reuters {
	return &Reuters{
		oldArticleTitles: oldArticleTitles,
		country:          country,
	}
}

func GetCountry(code uint8) Country {
	switch code {
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
