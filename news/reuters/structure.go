package reuters

type ReutersCategory struct {
	Type string `json:"type"`
	Data struct {
		Stories []ReutersArticleMetadata `json:"stories"`
	} `json:"data"`
}

type ReutersWebCategory struct {
	Result struct {
		Articles []ReutersArticleMetadata `json:"articles"`
	} `json:"result"`
}

type ReutersArticleMetadata struct {
	Title        string `json:"title"`
	URL          string `json:"url"`
	SectionURL   string `json:"section_url"`
	CanonicalURL string `json:"canonical_url"`
}

type ReutersMobileArticle struct {
	Type string `json:"type"`
	Data struct {
		ArticleData ReutersArticle `json:"article"`
	} `json:"data"`
}

type ReutersWebArticle struct {
	Result ReutersArticle `json:"result"`
}

type ReutersArticle struct {
	Title                string                  `json:"title"`
	ContentElements      []ReutersContentElement `json:"content_elements"`
	AdditionalProperties struct {
		ArticleProperties struct {
			Place string `json:"place"`
		} `json:"article_properties"`
	} `json:"additional_properties"`
	Thumbnail ReutersThumbnail `json:"thumbnail"`
}

type ReutersContentElement struct {
	Type    string `json:"type"`
	Content any    `json:"content"`
}

type ReutersThumbnail struct {
	ID      string `json:"id"`
	URL     string `json:"url"`
	Caption string `json:"caption"`
	AltText string `json:"alt_text"`
}
