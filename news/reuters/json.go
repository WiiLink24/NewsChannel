package reuters

import (
	"NewsChannel/news"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"slices"
	"strings"
)

func (r *Reuters) getMobileArticles(url string, topic news.Topic) ([]news.Article, error) {
	data, err := news.HttpGet(url, "ReutersNews/7.6.0 iPad8,6 iPadOS/18.1 CFNetwork/1.0 Darwin/24.1.0")
	if err != nil {
		return nil, err
	}

	var categoryData []ReutersCategory
	err = json.Unmarshal(data, &categoryData)
	if err != nil {
		return nil, err
	}

	var allowedTypes = []string{
		"story-cluster",
		"latest-stories",
	}

	// Iterate over the article block
	var articles []news.Article
	for _, v := range categoryData {
		if !slices.Contains(allowedTypes, v.Type) {
			continue
		}

		stories := v.Data.Stories
		for _, story := range stories {
			title := news.SanitizeText(story.Title)
			// Compare previous articles to see if we have a duplicate.
			if news.IsDuplicateArticle(r.oldArticleTitles, title) {
				return nil, nil
			}
			r.oldArticleTitles = append(r.oldArticleTitles, title)

			// Ignore podcasts
			if story.SectionURL == "/podcasts/" {
				return nil, nil
			}

			// The article is nested inside a "templates" list, with the data we require in the 1st index.
			// I (Noah) refer to this as bad because it returns the web page, rather than the mobile API page.
			// The mobile API is much easier to parse.
			articlePath := story.URL
			articleURL := fmt.Sprintf("https://www.reuters.com/mobile/v1%s", articlePath)
			articleData, err := news.HttpGet(articleURL, "ReutersNews/7.6.0 iPad8,6 iPadOS/18.1 CFNetwork/1.0 Darwin/24.1.0")
			if err != nil {
				return nil, err
			}

			// Parse article JSON
			var article []ReutersMobileArticle
			err = json.Unmarshal(articleData, &article)
			if err != nil {
				var serr *json.SyntaxError
				if errors.As(err, &serr) {
					return nil, nil
				}

				return nil, err
			}

			// Iterate until we find the "article_detail" key
			for _, child := range article {
				if child.Type != "article_detail" {
					continue
				}
				article, err := r.createArticle(child.Data.ArticleData, topic)
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
	}

	return articles, nil
}

func (r *Reuters) getWebArticles(url string, topic news.Topic) ([]news.Article, error) {
	data, err := news.HttpGet(url)
	if err != nil {
		return nil, err
	}

	var categoryData ReutersWebCategory
	err = json.Unmarshal(data, &categoryData)
	if err != nil {
		return nil, err
	}

	// Iterate over the article block
	var articles []news.Article
	stories := categoryData.Result.Articles
	for _, story := range stories {
		title := news.SanitizeText(story.Title)
		// Compare previous articles to see if we have a duplicate.
		if news.IsDuplicateArticle(r.oldArticleTitles, title) {
			return nil, nil
		}
		r.oldArticleTitles = append(r.oldArticleTitles, title)

		articlePath := story.CanonicalURL
		articleURL := fmt.Sprintf("https://jp.reuters.com/pf/api/v3/content/fetch/article-by-id-or-url-v1?query={\"website_url\":\"%s\",\"website\":\"reuters-japan\"}", articlePath)
		articleData, err := news.HttpGet(articleURL)
		if err != nil {
			return nil, err
		}

		// Parse article JSON
		var articleJSON ReutersWebArticle
		err = json.Unmarshal(articleData, &articleJSON)
		if err != nil {
			var serr *json.SyntaxError
			if errors.As(err, &serr) {
				return nil, nil
			}

			return nil, err
		}

		article, err := r.createArticle(articleJSON.Result, topic)
		if err != nil {
			return nil, err
		}
		if article == nil {
			continue
		}

		articles = append(articles, *article)
		return articles, nil
	}

	return articles, nil
}

func (r *Reuters) createArticle(story ReutersArticle, topic news.Topic) (*news.Article, error) {
	content, err := parseArticle(story.ContentElements)
	if err != nil {
		return nil, err
	}

	// Possible there is no text?
	if len(content) == 0 {
		return nil, nil
	}

	location, err := r.getLocation(story.AdditionalProperties.ArticleProperties.Place)
	if err != nil {
		return nil, err
	}

	// Finally get the thumbnail.
	thumbnail, err := getThumbnail(story.Thumbnail)
	if err != nil {
		return nil, err
	}

	return &news.Article{
		Title:     news.SanitizeText(story.Title),
		Content:   &content,
		Topic:     topic,
		Location:  location,
		Thumbnail: thumbnail,
	}, nil
}

func parseArticle(elements []ReutersContentElement) (string, error) {
	var ret string
	for _, content := range elements {
		if content.Type != "paragraph" {
			continue
		}

		// Sanitize paragraph
		unSanitized := content.Content.(string)
		sanitized := news.SanitizeText(unSanitized)

		ret += sanitized
		ret += "\n\n"
	}

	ret = strings.TrimSpace(ret)

	return ret, nil
}

func getThumbnail(thumbnail ReutersThumbnail) (*news.Thumbnail, error) {
	// Don't add Reuters logo as image
	if thumbnail.ID == "466BJJQ7PVGY5O53NZ3KL65MHM" || len(thumbnail.URL) == 0 {
		return nil, nil
	}

	data, err := news.HttpGet(thumbnail.URL, "ReutersNews/7.6.0 iPad8,6 iPadOS/18.1 CFNetwork/1.0 Darwin/24.1.0")
	if err != nil {
		return nil, err
	}

	if len(data) == 0 {
		return nil, nil
	}

	var caption string
	if len(thumbnail.Caption) != 0 {
		caption = thumbnail.Caption
	} else if len(thumbnail.AltText) != 0 {
		caption = thumbnail.AltText
	}

	return &news.Thumbnail{
		Image:   news.ConvertImage(data),
		Caption: news.SanitizeText(caption),
	}, nil
}

func (r *Reuters) getLocation(location string) (*news.Location, error) {
	if r.language == English {
		if strings.Contains(location, "Reuters") || len(location) == 0 {
			return nil, nil
		}
		locations := strings.Split(location, "/")

		// Use the new dynamic location function that includes OSM API fallback
		return news.GetLocationForExtractedLocation(locations, "en"), nil
	} else {
		datelineRegex := regexp.MustCompile(`([\[|［])(.*?)[０-９]`)
		locationString := datelineRegex.FindStringSubmatch(location)
		if len(locationString) > 2 && len(locationString[2]) > 0 {
			splitter := func(r rune) bool {
				return r == '/' || r == '／'
			}

			locationList := strings.FieldsFunc(locationString[2], splitter)
			return news.GetLocationForExtractedLocation(locationList, "jp"), nil
		}
	}
	return nil, nil
}
