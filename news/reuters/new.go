package reuters

import "NewsChannel/news"

// Country represents a country in the Reuters API
type Country string

const (
	UnitedStates = "us"
)

type Reuters struct {
	country Country
	news.Source
}

func NewReuters(country Country) *Reuters {
	return &Reuters{
		country: country,
	}
}
