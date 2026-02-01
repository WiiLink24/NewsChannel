package ap

import (
	"NewsChannel/news"
	"encoding/xml"
	"errors"
	"log"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// RSS structures for parsing AP XML feeds
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
	GUID        string `xml:"guid"`
}

func (a *AP) getArticles(url string, topic news.Topic) ([]news.Article, error) {
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
	for _, item := range rss.Channel.Items {
		title := news.SanitizeText(item.Title)
		// Check for duplicates
		if news.IsDuplicateArticle(a.oldArticleTitles, title) {
			continue
		}

		// Get full article content by scraping the link
		content, location, thumbnail, err := a.getFullArticle(item.Link)
		if err != nil {
			return nil, err
		}

		if *content == "" {
			continue
		}

		article := news.Article{
			Title:     title,
			Content:   content,
			Topic:     topic,
			Location:  location,
			Thumbnail: thumbnail,
		}

		articles = append(articles, article)
		break
	}

	return articles, nil
}

func (a *AP) getFullArticle(articleURL string) (*string, *news.Location, *news.Thumbnail, error) {
	if articleURL == "" {
		return nil, nil, nil, errors.New("empty articleURL")
	}

	data, err := news.HttpGet(articleURL)
	if err != nil {
		return nil, nil, nil, err
	}

	html := string(data)

	content, locationString, err := a.extractArticleBody(html)
	if err != nil {
		return nil, nil, nil, err
	}

	var location *news.Location
	if locationString != nil {
		location = news.GetLocationForExtractedLocation([]string{*locationString}, "en")
	}

	thumbnail := a.extractThumbnail(html)

	return content, location, thumbnail, nil
}

func (a *AP) extractArticleBody(html string) (*string, *string, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return nil, nil, err
	}

	var contentSlice []string

	// Select the main article body div
	doc.Find(`div.RichTextStoryBody`).Each(func(i int, section *goquery.Selection) {
		section.Find("p").EachWithBreak(func(j int, elem *goquery.Selection) bool {
			text := strings.TrimSpace(elem.Text())
			if text == "___" {
				// We have reached the article footer
				return false
			}

			if text != "" {
				contentSlice = append(contentSlice, text)
			}
			return true
		})
	})

	var content string
	for _, paragraph := range contentSlice {
		content += news.SanitizeText(paragraph)
		content += "\n\n"
	}

	// Get the location
	locationRegex := regexp.MustCompile(`(.*?) \(AP\) — `)
	location := locationRegex.FindStringSubmatch(contentSlice[0])
	if len(location) > 1 && len(location[1]) > 0 {
		return &content, &location[1], nil
	}

	return &content, nil, nil
}

func (a *AP) extractThumbnail(html string) *news.Thumbnail {
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

	imageData, err := news.HttpGet(imageURL)
	if err != nil || len(imageData) == 0 {
		return nil
	}

	caption := ""

	doc.Find(`meta[property="og:image:alt"]`).Each(func(i int, s *goquery.Selection) {
		if content, exists := s.Attr("content"); exists {
			caption = strings.TrimSpace(content)
			return
		}
	})

	return &news.Thumbnail{
		Image:   news.ConvertImage(imageData),
		Caption: news.SanitizeText(caption),
	}
}
