package nos

import (
	"NewsChannel/news"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"net/http"
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

func httpGet(url string) ([]byte, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil || resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP request failed: %v", err)
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func (f *nos) getArticles(url string, topic news.Topic) ([]news.Article, error) {
	// Fetch RSS XML
	data, err := httpGet(url)
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

		// Check for duplicates
		if news.IsDuplicateArticle(f.oldArticleTitles, item.Title) {
			continue
		}

		// Extract content from RSS
		content := f.extractContentFromDescription(item.Description)

		// Get location by scraping the article page for meta tags
		location := f.getLocationFromArticlePage(item.Link)

		// Get thumbnail from RSS
		var thumbnail *news.Thumbnail
		if item.Enclosure.URL != "" && strings.Contains(item.Enclosure.Type, "image") {
			imageData, err := httpGet(item.Enclosure.URL)
			if err == nil && len(imageData) > 0 {
				thumbnail = &news.Thumbnail{
					Image:   news.ConvertImage(imageData),
					Caption: "",
				}
			}
		}

		// Skip if no content
		if len(strings.TrimSpace(content)) == 0 {
			continue
		}

		article := news.Article{
			Title:     strings.TrimSpace(item.Title),
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
	doc, _ := goquery.NewDocumentFromReader(strings.NewReader(description))

	var builder strings.Builder

	text := strings.TrimSpace(doc.Text())
	if text != "" {
		builder.WriteString(text)
	}

	content := strings.TrimSpace(builder.String())
	return content
}

func (f *nos) getLocationFromArticlePage(articleURL string) *news.Location {
	if articleURL == "" {
		return nil
	}

	data, err := httpGet(articleURL)
	if err != nil {
		log.Printf("Failed to fetch article page for location: %v", err)
		return nil
	}

	html := string(data)
	return f.extractLocationFromContent(html)
}

func (f *nos) extractLocationFromContent(html string) *news.Location {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		log.Println("Failed to parse HTML:", err)
		return nil
	}

	// Try to find location in meta keywords
	var foundLocation *news.Location
	doc.Find(`meta[name="keywords"]`).EachWithBreak(func(i int, s *goquery.Selection) bool {
		if content, exists := s.Attr("content"); exists {
			keywords := strings.Split(content, ",")
			for _, keyword := range keywords {
				keyword = strings.TrimSpace(keyword)
				if keyword != "" {
					// Try to find location match
					if loc := news.GetLocationForExtractedLocation(keyword); loc != nil {
						foundLocation = loc
						return false
					}
				}
			}
		}
		return true
	})

	return foundLocation
}
