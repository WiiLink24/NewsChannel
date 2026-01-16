package france24

import (
	"NewsChannel/news"
	"encoding/xml"
	"log"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

var htmlRegex = regexp.MustCompile("<.*?>")

// RSS structures for parsing France24 XML feeds
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
	Title          string         `xml:"title"`
	Description    string         `xml:"description"`
	Link           string         `xml:"link"`
	PubDate        string         `xml:"pubDate"`
	GUID           string         `xml:"guid"`
	Category       string         `xml:"category"`
	Creator        string         `xml:"creator"`
	MediaThumbnail MediaThumbnail `xml:"thumbnail"`
}

type MediaThumbnail struct {
	URL string `xml:"url,attr"`
}

type france24 struct {
	oldArticleTitles []string
}

func (a *france24) getArticles(url string, topic news.Topic) ([]news.Article, error) {
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
		content = htmlRegex.ReplaceAllString(content, "")

		// Skip if no content
		if len(strings.TrimSpace(content)) == 0 {
			continue
		}

		// Use media thumbnail if available
		if thumbnail == nil && item.MediaThumbnail.URL != "" {
			imageData, err := news.HttpGet(item.MediaThumbnail.URL)
			if err == nil && len(imageData) > 0 {
				thumbnail = &news.Thumbnail{
					Image:   news.ConvertImage(imageData),
					Caption: "",
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

func (a *france24) getFullArticle(articleURL string) (content string, location *news.Location, thumbnail *news.Thumbnail) {
	if articleURL == "" {
		return "", nil, nil
	}

	data, err := news.HttpGet(articleURL)
	if err != nil {
		log.Printf("Failed to fetch article content from %s: %v", articleURL, err)
		return "", nil, nil
	}

	html := string(data)

	content = a.extractArticleBody(html)
	location = a.extractLocationFromContent(html)
	thumbnail = a.extractThumbnail(html)

	return content, location, thumbnail
}

func (a *france24) extractArticleBody(html string) string {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		log.Println("Failed to parse HTML:", err)
		return ""
	}

	var builder strings.Builder

	// Find all <p> tags with more than 50 characters and extract them (might not be the best approach)
	doc.Find("p").Each(func(i int, s *goquery.Selection) {
		text := strings.TrimSpace(s.Text())
		if text != "" && len(text) > 50 {
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

	// Remove unwated phrase at the end of the content
	content = strings.ReplaceAll(content, "Le contenu auquel vous tentez d'accéder n'existe pas ou n'est plus disponible.", "")

	return content
}

func (a *france24) extractLocationFromContent(html string) *news.Location {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		log.Println("Failed to parse HTML:", err)
		return nil
	}

	locationCandidates := make(map[string]bool)

	// Check meta candidates for location
	doc.Find(`meta[name="news_keywords"]`).Each(func(i int, s *goquery.Selection) {
		if content, exists := s.Attr("content"); exists {
			keywords := strings.Split(content, ",")
			for _, keyword := range keywords {
				keyword = strings.TrimSpace(keyword)
				if keyword != "" {
					locationCandidates[keyword] = true
				}
			}
		}
	})

	doc.Find(`meta[property="article:tag"]`).Each(func(i int, s *goquery.Selection) {
		if content, exists := s.Attr("content"); exists {
			content = strings.TrimSpace(content)
			if content != "" {
				locationCandidates[content] = true
			}
		}
	})

	// Try each candidate location
	var candidates []string
	for candidate := range locationCandidates {
		candidates = append(candidates, candidate)
	}

	if loc := news.GetLocationForExtractedLocation(candidates, "fr"); loc != nil {
		return loc
	}

	return nil
}

func (a *france24) extractThumbnail(html string) *news.Thumbnail {
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

	// Make sure URL is absolute
	if strings.HasPrefix(imageURL, "//") {
		imageURL = "https:" + imageURL
	} else if strings.HasPrefix(imageURL, "/") {
		imageURL = "https://www.france24.com" + imageURL
	}

	imageData, err := news.HttpGet(imageURL)
	if err != nil || len(imageData) == 0 {
		return nil
	}

	caption := ""
	doc.Find("figure.m-item-image figcaption.a-figcaption").Each(func(i int, s *goquery.Selection) {
		var captionParts []string
		s.Find("span").Each(func(j int, span *goquery.Selection) {
			text := strings.TrimSpace(span.Text())
			if text != "" {
				captionParts = append(captionParts, text)
			}
		})
		if len(captionParts) > 0 {
			caption = strings.Join(captionParts, " ")
			return
		}
	})

	return &news.Thumbnail{
		Image:   news.ConvertImage(imageData),
		Caption: caption,
	}
}
