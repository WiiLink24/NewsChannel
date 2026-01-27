package reuters

import (
	"NewsChannel/news"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strings"
)

func (r *Reuters) getArticles(url string, topic news.Topic) ([]news.Article, error) {
	data, err := news.HttpGet(url, "ReutersNews/7.6.0 iPad8,6 iPadOS/18.1 CFNetwork/1.0 Darwin/24.1.0")
	if err != nil {
		return nil, err
	}

	var root []map[string]any
	err = json.Unmarshal(data, &root)
	if err != nil {
		return nil, err
	}

	// Iterate over the article block
	var articles []news.Article
	for _, v := range root {
		if v["type"].(string) != "story-cluster" {
			continue
		}

		stories := v["data"].(map[string]any)["stories"].([]any)
		for _, story := range stories {
			article, err := r.createArticle(story.(map[string]any), topic)
			if err != nil {
				return nil, err
			}
			if article == nil {
				continue
			}

			articles = append(articles, *article)
			return articles, nil
		}
	}

	// If we make it here, all clusters are duplicates. Move onto latest articles.
	for _, v := range root {
		if v["type"].(string) != "latest-stories" {
			continue
		}

		stories := v["data"].(map[string]any)["stories"].([]any)
		for _, story := range stories {
			article, err := r.createArticle(story.(map[string]any), topic)
			if err != nil {
				return nil, err
			}
			if article == nil {
				continue
			}

			articles = append(articles, *article)
			return articles, nil
		}
	}

	return articles, nil
}

func (r *Reuters) createArticle(story map[string]any, topic news.Topic) (*news.Article, error) {
	title := news.SanitizeText(story["title"].(string))
	// Compare previous articles to see if we have a duplicate.
	if news.IsDuplicateArticle(r.oldArticleTitles, title) {
		return nil, nil
	}

	// Ignore podcasts
	if story["section_url"] == "/podcasts/" {
		return nil, nil
	}

	// The article is nested inside a "templates" list, with the data we require in the 1st index.
	// I (Noah) refer to this as bad because it returns the web page, rather than the mobile API page.
	// The mobile API is much easier to parse.
	articlePath := story["url"]
	articleURL := fmt.Sprintf("https://www.reuters.com/mobile/v1%s", articlePath)
	articleData, err := news.HttpGet(articleURL, "ReutersNews/7.6.0 iPad8,6 iPadOS/18.1 CFNetwork/1.0 Darwin/24.1.0")
	if err != nil {
		return nil, err
	}

	// Parse article JSON
	var articleJSON []map[string]any
	err = json.Unmarshal(articleData, &articleJSON)
	if err != nil {
		var serr *json.SyntaxError
		if errors.As(err, &serr) {
			return nil, nil
		}

		return nil, err
	}

	content, err := parseArticle(articleJSON)
	if err != nil {
		return nil, err
	}

	// Possible there is no text?
	if len(*content) == 0 {
		return nil, nil
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

	return &news.Article{
		Title:     title,
		Content:   content,
		Topic:     topic,
		Location:  location,
		Thumbnail: thumbnail,
	}, nil
}

func parseArticle(root []map[string]any) (*string, error) {
	// Iterate until we find the "article_detail" key
	var ret string
	for _, child := range root {
		if child["type"].(string) != "article_detail" {
			continue
		}

		for _, content := range child["data"].(map[string]any)["article"].(map[string]any)["content_elements"].([]any) {
			if content.(map[string]any)["type"].(string) != "paragraph" {
				continue
			}

			// Sanitize paragraph
			unSanitized := content.(map[string]any)["content"].(string)
			sanitized := news.SanitizeText(unSanitized)

			ret += sanitized
			ret += "\n\n"
		}
	}

	return &ret, nil
}

func getThumbnail(root []map[string]any) (*news.Thumbnail, error) {
	for _, child := range root {
		if child["type"].(string) != "article_detail" {
			continue
		}

		// Don't add Reuters logo as image
		if child["data"].(map[string]any)["article"].(map[string]any)["thumbnail"].(map[string]any)["id"] != nil {
			if child["data"].(map[string]any)["article"].(map[string]any)["thumbnail"].(map[string]any)["id"].(string) == "466BJJQ7PVGY5O53NZ3KL65MHM" {
				return nil, nil
			}
		}

		thumbnailURL := child["data"].(map[string]any)["article"].(map[string]any)["thumbnail"].(map[string]any)["resizer_url"].(string)

		// Add the required params
		parsedURL, err := url.Parse(thumbnailURL)
		if err != nil {
			return nil, err
		}
		query := parsedURL.Query()
		query.Add("width", "200")
		query.Add("height", "200")
		parsedURL.RawQuery = query.Encode()
		thumbnailURL = parsedURL.String()

		data, err := news.HttpGet(thumbnailURL, "ReutersNews/7.6.0 iPad8,6 iPadOS/18.1 CFNetwork/1.0 Darwin/24.1.0")
		if err != nil {
			return nil, err
		}

		if len(data) == 0 {
			return nil, nil
		}

		caption := ""
		if child["data"].(map[string]any)["article"].(map[string]any)["thumbnail"].(map[string]any)["caption"] != nil {
			caption = child["data"].(map[string]any)["article"].(map[string]any)["thumbnail"].(map[string]any)["caption"].(string)
		}

		return &news.Thumbnail{
			Image:   news.ConvertImage(data),
			Caption: caption,
		}, nil
	}

	return nil, nil
}

func getLocation(root []map[string]any) (*news.Location, error) {
	for _, child := range root {
		if child["type"].(string) != "article_detail" {
			continue
		}

		if child["data"].(map[string]any)["article"].(map[string]any)["additional_properties"].(map[string]any)["article_properties"].(map[string]any)["place"] == nil {
			return nil, nil
		}

		location := child["data"].(map[string]any)["article"].(map[string]any)["additional_properties"].(map[string]any)["article_properties"].(map[string]any)["place"].(string)

		// Sometimes the place property is just the date, followed by "(Reuters)"
		if strings.Contains(location, "(Reuters)") {
			return nil, nil
		}

		locations := strings.Split(location, "/")
		// Use the new dynamic location function that includes OSM API fallback
		return news.GetLocationForExtractedLocation(locations, "en"), nil
	}

	return nil, nil
}
