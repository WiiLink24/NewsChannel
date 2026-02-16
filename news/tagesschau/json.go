package tagesschau

import (
	"NewsChannel/news"
	"encoding/json"
	"errors"
)

func (r *Tagesschau) getArticles(url string, topic news.Topic, storyKey string) ([]news.Article, error) {
	data, err := news.HttpGet(url)
	if err != nil {
		return nil, err
	}

	var root map[string]any
	err = json.Unmarshal(data, &root)
	if err != nil {
		return nil, err
	}

	stories := root[storyKey].([]any)

	// Iterate over the article block
	var articles []news.Article
	for _, story := range stories {
		title := news.SanitizeText(story.(map[string]any)["title"].(string))
		// Compare previous articles to see if we have a duplicate.
		if news.IsDuplicateArticle(r.oldArticleTitles, title) {
			continue
		}
		r.oldArticleTitles = append(r.oldArticleTitles, title)

		// Ignore non-articles
		if story.(map[string]any)["type"].(string) != "story" {
			continue
		}

		articleURL := story.(map[string]any)["details"].(string)
		articleData, err := news.HttpGet(articleURL)
		if err != nil {
			return nil, err
		}

		// Parse article JSON
		var articleJSON map[string]any
		err = json.Unmarshal(articleData, &articleJSON)
		if err != nil {
			var serr *json.SyntaxError
			if errors.As(err, &serr) {
				continue
			}

			return nil, err
		}

		content, err := parseArticle(articleJSON)
		if err != nil {
			return nil, err
		}

		// Possible there is no text?
		if len(*content) == 0 {
			continue
		}

		location, err := getLocation(articleJSON)
		if err != nil {
			return nil, err
		}

		// Finally get the thumbnail.
		thumbnail, err := getThumbnail(articleJSON)
		if err != nil {
			return nil, err
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

func parseArticle(root map[string]any) (*string, error) {
	// Iterate through text content
	var ret string
	for _, content := range root["content"].([]any) {
		if !allowedTypes[content.(map[string]any)["type"].(string)] {
			continue
		}

		var unSanitized string
		if content.(map[string]any)["type"].(string) == "quotation" {
			unSanitized = content.(map[string]any)["quotation"].(map[string]any)["text"].(string)
		} else {
			unSanitized = content.(map[string]any)["value"].(string)
		}

		sanitized := news.SanitizeText(unSanitized)

		ret += sanitized
		ret += "\n\n"
	}

	return &ret, nil
}

func getThumbnail(root map[string]any) (*news.Thumbnail, error) {
	if root["teaserImage"] == nil {
		return nil, nil
	}

	image := root["teaserImage"].(map[string]any)
	if image["imageVariants"] == nil {
		return nil, nil
	}

	// Ignore Tagesschau logo
	if image["alttext"] != nil {
		if image["alttext"].(string) == "Globus auf blauem Hintergrund mit tagesschau-Schriftzug" {
			return nil, nil
		}
	}

	acceptedThumbnails := []string{
		"1x1-840", "1x1-640", "1x1-432", "1x1-256", "1x1-144",
	}
	var thumbnailURL string

	// Get highest res 1x1 ratio image URL
	for _, thumbnail := range acceptedThumbnails {
		if image["imageVariants"].(map[string]any)[thumbnail] != nil {
			thumbnailURL = image["imageVariants"].(map[string]any)[thumbnail].(string)
			break
		}
	}

	if thumbnailURL == "" {
		return nil, nil
	}

	data, err := news.HttpGet(thumbnailURL)
	if err != nil {
		return nil, err
	}

	if len(data) == 0 {
		return nil, nil
	}

	caption := ""
	if image["title"] != nil {
		caption = image["title"].(string)
	} else if image["alttext"] != nil {
		caption = image["alttext"].(string)
	}

	return &news.Thumbnail{
		Image:   news.ConvertImage(data),
		Caption: news.SanitizeText(caption),
	}, nil
}

func getLocation(root map[string]any) (*news.Location, error) {
	var tags []string
	for _, tag := range root["tags"].([]any) {
		tags = append(tags, tag.(map[string]any)["tag"].(string))
	}

	if len(tags) != 0 {
		return news.GetLocationForExtractedLocation(tags, "de"), nil
	}

	return nil, nil
}

var allowedTypes = map[string]bool{
	"text":      true,
	"headline":  true,
	"quotation": true,
}
