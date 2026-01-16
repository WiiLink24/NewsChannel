package nhk

import (
	"NewsChannel/news"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

var htmlRegex = regexp.MustCompile("<.*?>")

// RSS structures for parsing NHK XML feeds
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
}

type nhk struct {
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

func (f *nhk) getArticles(url string, topic news.Topic) ([]news.Article, error) {
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

		title := strings.TrimSpace(item.Title)
		// Check for duplicates
		if news.IsDuplicateArticle(f.oldArticleTitles, title) {
			continue
		}

		// Get full article content by scraping the link
		content, location, thumbnail := f.getFullArticle(item.Link)

		// Use description as fallback if content fetch fails
		if content == "" {
			content = item.Description
		}

		// Clean HTML tags from content
		content = htmlRegex.ReplaceAllString(content, "")

		// Skip if no content
		if len(strings.TrimSpace(content)) == 0 {
			continue
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

func (f *nhk) getFullArticle(articleURL string) (content string, location *news.Location, thumbnail *news.Thumbnail) {
	if articleURL == "" {
		return "", nil, nil
	}

	data, err := httpGet(articleURL)
	if err != nil {
		log.Printf("Failed to fetch article content from %s: %v", articleURL, err)
		return "", nil, nil
	}

	html := string(data)

	content = f.extractArticleBody(html)
	location = f.extractLocationFromContent(html)
	thumbnail = f.extractThumbnail(html)

	return content, location, thumbnail
}

func (f *nhk) extractArticleBody(html string) string {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		log.Println("Failed to parse HTML:", err)
		return ""
	}

	var builder strings.Builder

	// Summary
	doc.Find("p.content--summary").Each(func(i int, s *goquery.Selection) {
		text := strings.TrimSpace(s.Text())
		if text != "" {
			builder.WriteString(text)
			builder.WriteString("\n\n")
		}
	})

	// Main article body
	doc.Find("section.content--body .body-text p").Each(func(i int, s *goquery.Selection) {
		text := strings.TrimSpace(s.Text())
		if text != "" {
			builder.WriteString(text)
			builder.WriteString("\n\n")
		}
	})

	content := strings.TrimSpace(builder.String())

	// Clean up HTML entities
	content = strings.ReplaceAll(content, "&nbsp;", " ")
	content = strings.ReplaceAll(content, "&amp;", "&")
	content = strings.ReplaceAll(content, "&lt;", "<")
	content = strings.ReplaceAll(content, "&gt;", ">")
	content = strings.ReplaceAll(content, "&quot;", "\"")
	content = strings.ReplaceAll(content, "&#39;", "'")
	content = strings.ReplaceAll(content, "&rsquo;", "'")
	content = strings.ReplaceAll(content, "&lsquo;", "'")

	return content
}

func (f *nhk) extractLocationFromContent(html string) *news.Location {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		log.Println("Failed to parse HTML:", err)
		return nil
	}

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

	if loc := news.GetLocationForExtractedLocation(candidates, "ja"); loc != nil {
		return loc
	}

	return nil
}

func (f *nhk) extractThumbnail(html string) *news.Thumbnail {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		log.Println("Failed to parse HTML:", err)
		return nil
	}

	var imageURL string

	// Get image using og:image meta tag
	if imageURL == "" {
		doc.Find(`meta[property="og:image"]`).EachWithBreak(func(i int, s *goquery.Selection) bool {
			if content, exists := s.Attr("content"); exists && strings.TrimSpace(content) != "" {
				imageURL = content
				return false
			}
			return true
		})
	}

	if imageURL == "" {
		return nil
	}

	// Make sure URL is absolute
	if strings.HasPrefix(imageURL, "//") {
		imageURL = "https:" + imageURL
	} else if strings.HasPrefix(imageURL, "/") {
		imageURL = "https://www3.nhk.or.jp" + imageURL
	}

	imageData, err := httpGet(imageURL)
	if err != nil || len(imageData) == 0 {
		return nil
	}

	return &news.Thumbnail{
		Image:   news.ConvertImage(imageData),
		Caption: "",
	}
}
