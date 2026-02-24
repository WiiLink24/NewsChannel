package nos

import (
	"NewsChannel/news"
	"encoding/xml"
	"log"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// RSS structures for parsing NOS XML feeds
type RSS struct {
	XMLName xml.Name `xml:"rss"`
	Channel Channel  `xml:"channel"`
}

type Channel struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	Items       []Item `xml:"item"`
}

type Item struct {
	Title       string `xml:"title"`
	Description string `xml:"description"`
	Link        string `xml:"link"`
	PubDate     string `xml:"pubDate"`
	GUID        string `xml:"guid"`
	Enclosure   struct {
		URL  string `xml:"url,attr"`
		Type string `xml:"type,attr"`
	} `xml:"enclosure"`
}

type nos struct {
	oldArticleTitles []string
}

func (f *nos) getArticles(url string, topic news.Topic) ([]news.Article, error) {
	// Fetch RSS XML
	data, err := news.HttpGet(url)
	if err != nil {
		return nil, err
	}

	// Parse RSS XML
	var rss RSS
	err = xml.Unmarshal(data, &rss)
	if err != nil {
		return nil, err
	}

	var articles []news.Article
	for i, item := range rss.Channel.Items {
		if i >= 1 {
			break
		}

		title := news.SanitizeText(item.Title)
		// Check for duplicates
		if news.IsDuplicateArticle(f.oldArticleTitles, title) {
			continue
		}
		f.oldArticleTitles = append(f.oldArticleTitles, title)

		// Extract content from RSS
		content := f.extractContentFromDescription(item.Description)

		// Skip if no content
		if len(content) == 0 {
			continue
		}

		// Get location by scraping the article page for meta tags
		location, err := f.getLocationFromArticlePage(content, item.Link)
		if err != nil {
			log.Printf("Failed to get location for %s: %s", item.Link, err)
		}

		// Get thumbnail from RSS
		var thumbnail *news.Thumbnail
		if item.Enclosure.URL != "" && strings.Contains(item.Enclosure.Type, "image") {
			imageData, err := news.HttpGet(item.Enclosure.URL)
			if err == nil && len(imageData) > 0 {
				caption := f.extractImageCaption(item.Link)
				thumbnail = &news.Thumbnail{
					Image:   news.ConvertImage(imageData),
					Caption: news.SanitizeText(caption),
				}
			}
		}

		article := news.Article{
			Title:     title,
			Content:   &content,
			Topic:     topic,
			Location:  location,
			Thumbnail: thumbnail,
		}

		articles = append(articles, article)
	}

	return articles, nil
}

func (f *nos) extractContentFromDescription(description string) string {
	var ret string

	doc, _ := goquery.NewDocumentFromReader(strings.NewReader(description))
	text := strings.Split(doc.Text(), "\n")

	for _, paragraph := range text {
		sanitized := news.SanitizeText(paragraph)
		ret += sanitized
		ret += "\n\n"
	}

	return strings.TrimSpace(ret)
}

func (f *nos) getLocationFromArticlePage(content string, articleURL string) (*news.Location, error) {
	// First check if Google Maps is enabled.
	// It returns far better locations for non-English languages.
	if news.UseGmaps {
		return news.GetGmapsLocation(content, "nl"), nil
	}

	if articleURL == "" {
		return nil, nil
	}

	data, err := news.HttpGet(articleURL)
	if err != nil {
		return nil, err
	}

	html := string(data)
	return f.extractLocationFromContent(html)
}

func (f *nos) extractLocationFromContent(html string) (*news.Location, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return nil, err
	}

	// Try to find location in meta keywords
	var candidates []string
	doc.Find(`meta[name="keywords"]`).EachWithBreak(func(i int, s *goquery.Selection) bool {
		if content, exists := s.Attr("content"); exists {
			keywords := strings.Split(content, ",")
			for _, keyword := range keywords {
				keyword = strings.TrimSpace(keyword)
				if keyword != "" {
					candidates = append(candidates, keyword)
				}
			}
		}
		return true
	})

	return news.GetLocationForExtractedLocation(candidates, "nl"), nil
}

func (f *nos) extractImageCaption(articleURL string) string {
	if articleURL == "" {
		return ""
	}

	data, err := news.HttpGet(articleURL)
	if err != nil {
		return ""
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(data)))
	if err != nil {
		return ""
	}

	ariaLabel, exists := doc.Find("button[aria-label*='copyright']").Attr("aria-label")
	if exists && ariaLabel != "" {
		if strings.Contains(ariaLabel, "copyright:") {
			parts := strings.Split(ariaLabel, "copyright:")
			if len(parts) > 1 {
				return strings.TrimSpace(parts[1])
			}
		}
	}

	return ""
}
