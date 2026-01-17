package ansa

import (
	"NewsChannel/news"
	"encoding/xml"
	"log"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// RSS structures for parsing ANSA XML feeds
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

func (a *ANSA) getArticles(url string, topic news.Topic) ([]news.Article, error) {
	// Fetch RSS XML
	data, err := news.HttpGet(url, "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")
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
		if news.IsDuplicateArticle(a.oldArticleTitles, title) {
			continue
		}

		// Get full article content by scraping the link
		content, location, thumbnail := a.getFullArticle(item.Link)

		// Use description as fallback if content fetch fails
		if content == "" {
			content = item.Description
		}

		// Clean HTML tags from content
		content = news.SanitizeText(content)

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

func (a *ANSA) getFullArticle(articleURL string) (content string, location *news.Location, thumbnail *news.Thumbnail) {
	if articleURL == "" {
		return "", nil, nil
	}

	data, err := news.HttpGet(articleURL, "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")
	if err != nil {
		log.Printf("Failed to fetch article content from %s: %v", articleURL, err)
		return "", nil, nil
	}

	html := string(data)

	content = a.extractArticleBody(html)
	location = a.extractLocationFromTags(html)
	thumbnail = a.extractThumbnail(html)

	return content, location, thumbnail
}

func (a *ANSA) extractArticleBody(html string) string {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		log.Println("Failed to parse HTML:", err)
		return ""
	}

	var builder strings.Builder

	// Select the main article body div
	doc.Find(`div.post-single-text.rich-text.news-txt[itemprop="articleBody"]`).Each(func(i int, section *goquery.Selection) {
		section.Find(`div.rich-text, #piano-container`).Remove()

		section.Find("p, div, span, h1, h2, h3, h4, h5, h6").Each(func(j int, elem *goquery.Selection) {
			// Split on tabs (ANSA add a tab at the start of each paragraph)
			paragraphSplit := strings.Split(elem.Text(), "\t")

			for _, paragraph := range paragraphSplit {
				text := strings.TrimSpace(paragraph)
				text = strings.ReplaceAll(text, "\n", " ")
				if text != "" {
					builder.WriteString(text + "\n\n")
				}
			}
		})
	})

	result := strings.TrimSpace(builder.String())
	return result
}

func (a *ANSA) extractLocationFromTags(html string) *news.Location {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		log.Println("Failed to parse HTML:", err)
		return nil
	}

	seen := make(map[string]bool)

	// Fetch all possible tags
	doc.Find(`a.tag`).Each(func(i int, s *goquery.Selection) {
		text := strings.TrimSpace(s.Text())
		if text != "" {
			seen[text] = true
		}
	})

	scriptTagRegex := regexp.MustCompile(`displayTags\("([^"]+)",\s*"[^"]*"\);`)
	scriptMatches := scriptTagRegex.FindAllStringSubmatch(html, -1)
	for _, match := range scriptMatches {
		if len(match) > 1 {
			tags := strings.Split(match[1], ",")
			for _, tag := range tags {
				tag = strings.TrimSpace(tag)
				if tag != "" {
					seen[tag] = true
				}
			}
		}
	}

	// Try each candidate tag as a location
	var tags []string
	for tag := range seen {
		tags = append(tags, tag)
	}

	return news.GetLocationForExtractedLocation(tags, "it")
}

func (a *ANSA) extractThumbnail(html string) *news.Thumbnail {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		log.Println("Failed to parse HTML:", err)
		return nil
	}

	var imageURL string

	doc.Find(`meta[property="og:image"]`).EachWithBreak(func(i int, s *goquery.Selection) bool {
		if content, exists := s.Attr("content"); exists && strings.TrimSpace(content) != "" {
			imageURL = content
			return false
		}
		return true
	})

	if imageURL == "" {
		return nil
	}

	imageData, err := news.HttpGet(imageURL, "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")
	if err != nil || len(imageData) == 0 {
		return nil
	}

	caption := ""

	// Try to find caption in figure data-caption attribute
	doc.Find("figure.image a[data-caption]").Each(func(i int, s *goquery.Selection) {
		if dataCap, exists := s.Attr("data-caption"); exists {
			caption = strings.TrimSpace(dataCap)
			return
		}
	})

	// Fallback to .image-caption if data-caption not found
	if caption == "" {
		doc.Find(".image-caption").Each(func(i int, s *goquery.Selection) {
			caption = strings.TrimSpace(s.Text())
		})
	}

	return &news.Thumbnail{
		Image:   news.ConvertImage(imageData),
		Caption: caption,
	}
}
