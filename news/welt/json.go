package welt

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

// RSS structures for parsing Welt XML feeds
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
	Title         string       `xml:"title"`
	Description   string       `xml:"description"`
	Link          string       `xml:"link"`
	PubDate       string       `xml:"pubDate"`
	GUID          string       `xml:"guid"`
	Category      string       `xml:"category"`
	Creator       string       `xml:"creator"`
	MediaContent  MediaContent `xml:"content"`
	MediaKeywords string       `xml:"keywords"`
}

type MediaContent struct {
	URL       string         `xml:"url,attr"`
	Type      string         `xml:"type,attr"`
	Thumbnail MediaThumbnail `xml:"thumbnail"`
}

type MediaThumbnail struct {
	URL string `xml:"url,attr"`
}

type welt struct {
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

func (f *welt) getArticles(url string, topic news.Topic) ([]news.Article, error) {
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

		// Scrape the full article content
		content := f.scrapeArticleContent(item.Link)

		// Fallback to RSS description if scraping fails
		if content == "" {
			content = htmlRegex.ReplaceAllString(item.Description, "")
			content = f.cleanContent(content)
		}

		// Skip if no content
		if len(strings.TrimSpace(content)) == 0 {
			continue
		}

		var location *news.Location
		if item.MediaKeywords != "" {
			keywords := strings.Split(item.MediaKeywords, ",")
			for _, keyword := range keywords {
				keyword = strings.TrimSpace(keyword)
				if keyword != "" {
					location = news.GetLocationForExtractedLocation(keyword)
					if location != nil {
						break
					}
				}
			}
		}

		var thumbnail *news.Thumbnail
		if item.MediaContent.Thumbnail.URL != "" {
			imageData, err := httpGet(item.MediaContent.Thumbnail.URL)
			if err == nil && len(imageData) > 0 {
				convertedImage := news.ConvertImage(imageData)
				if convertedImage != nil {
					thumbnail = &news.Thumbnail{
						Image:   convertedImage,
						Caption: "",
					}
				}
			}
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

func (f *welt) scrapeArticleContent(articleURL string) string {
	if articleURL == "" {
		return ""
	}

	data, err := httpGet(articleURL)
	if err != nil {
		log.Printf("Failed to fetch article content from %s: %v", articleURL, err)
		return ""
	}

	html := string(data)
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		log.Printf("Failed to parse HTML from %s: %v", articleURL, err)
		return ""
	}

	var builder strings.Builder

	// Article intro
	doc.Find("div.c-article-page__intro p").Each(func(i int, s *goquery.Selection) {
		text := strings.TrimSpace(s.Text())
		if text != "" {
			builder.WriteString(text)
			builder.WriteString("\n\n")
		}
	})

	// Article content
	doc.Find("div.c-rich-text-renderer--article p").Each(func(i int, s *goquery.Selection) {
		text := strings.TrimSpace(s.Text())
		if text != "" && len(text) > 20 {
			builder.WriteString(text)
			builder.WriteString("\n\n")
		}
	})

	content := strings.TrimSpace(builder.String())
	return f.cleanContent(content)
}

func (f *welt) cleanContent(content string) string {
	content = strings.ReplaceAll(content, "&nbsp;", " ")
	content = strings.ReplaceAll(content, "&amp;", "&")
	content = strings.ReplaceAll(content, "&lt;", "<")
	content = strings.ReplaceAll(content, "&gt;", ">")
	content = strings.ReplaceAll(content, "&quot;", "\"")
	content = strings.ReplaceAll(content, "&#39;", "'")
	content = strings.ReplaceAll(content, "&rsquo;", "'")
	content = strings.ReplaceAll(content, "&lsquo;", "'")

	content = regexp.MustCompile(`\s+`).ReplaceAllString(content, " ")
	content = strings.TrimSpace(content)

	return content
}
