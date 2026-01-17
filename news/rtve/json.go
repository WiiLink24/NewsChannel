package rtve

import (
	"NewsChannel/news"
	"encoding/json"
	"strings"
)

// RTVEResponse represents the structure of RTVE API response
type RTVEResponse struct {
	Page struct {
		Items       []RTVEArticle `json:"items"`
		Number      int           `json:"number"`
		Size        int           `json:"size"`
		Offset      int           `json:"offset"`
		Total       int           `json:"total"`
		TotalPages  int           `json:"totalPages"`
		NumElements int           `json:"numElements"`
	} `json:"page"`
}

type RTVEArticle struct {
	ID                       string   `json:"id"`
	Title                    string   `json:"title"`
	AnteTitle                *string  `json:"anteTitle"`
	LongTitle                string   `json:"longTitle"`
	ShortTitle               *string  `json:"shortTitle"`
	Summary                  string   `json:"summary"`
	Text                     string   `json:"text"`
	Image                    string   `json:"image"`
	ImageSEO                 string   `json:"imageSEO"`
	HTMLUrl                  string   `json:"htmlUrl"`
	HTMLShortUrl             string   `json:"htmlShortUrl"`
	PublicationDate          string   `json:"publicationDate"`
	ModificationDate         string   `json:"modificationDate"`
	PublicationDateTimestamp int64    `json:"publicationDateTimestamp"`
	ContentType              string   `json:"contentType"`
	Language                 string   `json:"language"`
	MainCategory             string   `json:"mainCategory"`
	OtherTopicsName          []string `json:"otherTopicsName"`
	PubState                 struct {
		Code        string `json:"code"`
		Description string `json:"description"`
	} `json:"pubState"`
}

func (r *RTVE) getArticles(url string, topic news.Topic) ([]news.Article, error) {
	data, err := news.HttpGet(url)
	if err != nil {
		return nil, err
	}

	var response RTVEResponse
	err = json.Unmarshal(data, &response)
	if err != nil {
		return nil, err
	}

	var articles []news.Article
	for i, rtveArticle := range response.Page.Items {
		if i >= 1 { // Limit to 1 article per category like Reuters
			break
		}

		// Use longTitle if available, otherwise fallback to title
		title := rtveArticle.Title
		if rtveArticle.LongTitle != "" {
			title = rtveArticle.LongTitle
		}

		title = news.CleanHTMLEntities(title)

		// Check for duplicates
		if news.IsDuplicateArticle(r.oldArticleTitles, title) {
			continue
		}

		// Use the text field as content, clean HTML tags
		content := rtveArticle.Text
		if content == "" {
			// Fall back to summary if no text
			content = rtveArticle.Summary
		}

		// Cahnge "</p><p>" to "\n\n" for better readability
		content = strings.ReplaceAll(content, "</p><p>", "\n\n")

		content = news.CleanHTMLEntities(content)

		// Skip if no content
		if len(strings.TrimSpace(content)) == 0 {
			continue
		}

		// Get thumbnail - try imageSEO first, then image
		var thumbnail *news.Thumbnail
		imageURL := rtveArticle.ImageSEO
		if imageURL == "" {
			imageURL = rtveArticle.Image
		}

		if imageURL != "" {
			thumbnail, _ = r.getThumbnail(imageURL, rtveArticle.HTMLUrl)
		}

		// Parse location from content, category, and other topics
		location := r.extractLocation(content, rtveArticle.MainCategory, rtveArticle.OtherTopicsName)

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

func (r *RTVE) getThumbnail(imageURL string, articleURL string) (*news.Thumbnail, error) {
	if imageURL == "" {
		return nil, nil
	}

	// Ensure URL is absolute
	if !strings.HasPrefix(imageURL, "http") {
		if strings.HasPrefix(imageURL, "//") {
			imageURL = "https:" + imageURL
		} else if strings.HasPrefix(imageURL, "/") {
			imageURL = "https://img.rtve.es" + imageURL
		} else {
			imageURL = "https://img.rtve.es/" + imageURL
		}
	}

	data, err := news.HttpGet(imageURL)
	if err != nil {
		return nil, err
	}

	if len(data) == 0 {
		return nil, nil
	}

	caption := news.ExtractImageCaption(articleURL, "figcaption.figcaption span")

	return &news.Thumbnail{
		Image:   news.ConvertImage(data),
		Caption: caption,
	}, nil
}

func (r *RTVE) extractLocation(text, category string, otherTopics []string) *news.Location {
	// Extract location from a category path
	extractFromPath := func(path string) *news.Location {
		if path == "" {
			return nil
		}
		parts := strings.Split(path, "/")

		// Match the last part of the path recursively
		var candidates []string
		for i := len(parts) - 1; i >= 0; i-- {
			part := strings.TrimSpace(parts[i])
			if part != "" && part != "Noticias" && part != "Mundo" {
				if part != "Especiales" && part != "Nacional" && part != "Internacional" && part != "Tags Libres" {
					candidates = append(candidates, part)
				}
			}
		}

		return news.GetLocationForExtractedLocation(candidates, "es")
	}

	// Try to extract location from the main category
	if location := extractFromPath(category); location != nil {
		return location
	}

	// Main category doesn't contain a valid location, search through otherTopicsName
	for _, topic := range otherTopics {
		if location := extractFromPath(topic); location != nil {
			return location
		}
	}

	// No location found anywhere
	return nil
}
