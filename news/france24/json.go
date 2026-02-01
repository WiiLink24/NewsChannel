package france24

import (
	"NewsChannel/news"
	"encoding/xml"
	"log"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

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

		title := news.SanitizeText(item.Title)
		// Check for duplicates
		if news.IsDuplicateArticle(a.oldArticleTitles, title) {
			continue
		}
		a.oldArticleTitles = append(a.oldArticleTitles, title)

		// Get full article content by scraping the link
		content, location, thumbnail, err := a.getFullArticle(item.Link)
		if err != nil {
			return nil, err
		}

		contentString := *content

		// Use description as fallback if content fetch fails
		if contentString == "" {
			contentString = item.Description
		}

		// Clean HTML tags from content
		contentString = news.SanitizeText(contentString)

		// Skip if no content
		if len(strings.TrimSpace(contentString)) == 0 {
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
			Content:   &contentString,
			Topic:     topic,
			Location:  location,
			Thumbnail: thumbnail,
		}

		articles = append(articles, article)
	}

	return articles, nil
}

func (a *france24) getFullArticle(articleURL string) (*string, *news.Location, *news.Thumbnail, error) {
	if articleURL == "" {
		return nil, nil, nil, nil
	}

	data, err := news.HttpGet(articleURL)
	if err != nil {
		return nil, nil, nil, err
	}

	html := string(data)

	content, err := a.extractArticleBody(html)
	if err != nil {
		return nil, nil, nil, err
	}
	location := a.extractLocationFromContent(html)
	thumbnail := a.extractThumbnail(html)

	return content, location, thumbnail, nil
}

func (a *france24) extractArticleBody(html string) (*string, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return nil, err
	}

	var ret string

	// Find article body section
	doc.Find(`div.t-content__body`).Each(func(i int, section *goquery.Selection) {
		section.Find(`p.a-read-more`).Remove()
		section.Find("p").Each(func(j int, elem *goquery.Selection) {
			text := strings.TrimSpace(elem.Text())
			if text != "" {
				ret += text
				ret += "\n\n"
			}
		})
	})

	return &ret, nil
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

	return news.GetLocationForExtractedLocation(candidates, "fr")
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
		Caption: news.SanitizeText(caption),
	}
}
