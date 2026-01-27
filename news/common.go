package news

// Source represents a News source.
type Source interface {
	GetArticles() ([]Article, error)
	GetLogo() []byte
	GetCopyright() []uint16
}

type Article struct {
	Title     string
	Content   *string
	Topic     Topic
	Location  *Location
	Thumbnail *Thumbnail
}

type Thumbnail struct {
	Image   []byte
	Caption string
}

// Topic represents a news topic.
type Topic int

const (
	NationalNews Topic = iota
	InternationalNews
	Sports
	Entertainment
	Business
	Science
	Technology
)
