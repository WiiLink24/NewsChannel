package reutersjp

import (
	"NewsChannel/news"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func (r *ReutersJP) getArticles(url string, topic news.Topic) ([]news.Article, error) {
	data, err := news.HttpGet(url)
	if err != nil {
		return nil, err
	}

	var root map[string]any
	err = json.Unmarshal(data, &root)
	if err != nil {
		return nil, err
	}

	// Iterate over the article block
	var articles []news.Article
	stories := root["result"].(map[string]any)["articles"].([]any)
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

	return articles, nil
}

func (r *ReutersJP) createArticle(story map[string]any, topic news.Topic) (*news.Article, error) {
	title := news.SanitizeText(story["title"].(string))
	// Compare previous articles to see if we have a duplicate.
	if news.IsDuplicateArticle(r.oldArticleTitles, title) {
		return nil, nil
	}
	r.oldArticleTitles = append(r.oldArticleTitles, title)

	articlePath := story["canonical_url"]
	articleURL := fmt.Sprintf("https://jp.reuters.com%s", articlePath)
	articleData, err := news.HttpGet(articleURL)
	if err != nil {
		return nil, err
	}

	article := string(articleData)

	content, locationString, err := extractArticleBody(article)
	if err != nil {
		return nil, err
	}

	// Possible there is no text?
	if len(*content) == 0 {
		return nil, nil
	}

	var location *news.Location
	if locationString != nil {
		splitter := func(r rune) bool {
			return r == '/' || r == '／'
		}
		location = news.GetLocationForExtractedLocation(strings.FieldsFunc(*locationString, splitter), "jp")
	} else {
		location = nil
	}

	// Finally get the thumbnail.
	thumbnail, err := getThumbnail(story)
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

func extractArticleBody(html string) (*string, *string, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return nil, nil, err
	}

	var contentSlice []string

	// Select the main article body div
	doc.Find(`div.article-body-module__content__bnXL1`).Each(func(i int, section *goquery.Selection) {
		section.Find("div.article-body-module__paragraph__Ts-yF").Each(func(j int, elem *goquery.Selection) {
			text := strings.TrimSpace(elem.Text())
			if text != "" {
				contentSlice = append(contentSlice, text)
			}
		})
	})

	var result string

	for _, content := range contentSlice {
		result += news.SanitizeText(content)
		result += "\n\n"
	}

	// Get the location
	dateline := regexp.MustCompile(`([\[|［])(.*?)[０-９]`)
	location := dateline.FindStringSubmatch(contentSlice[0])
	if len(location) > 2 && len(location[2]) > 0 {
		return &result, &location[2], nil
	}

	return &result, nil, nil
}

func getThumbnail(story map[string]any) (*news.Thumbnail, error) {
	// Don't add Reuters logo as image
	if story["thumbnail"].(map[string]any)["id"] != nil {
		if story["thumbnail"].(map[string]any)["id"].(string) == "466BJJQ7PVGY5O53NZ3KL65MHM" {
			return nil, nil
		}
	}

	thumbnailURL := story["thumbnail"].(map[string]any)["url"].(string)

	data, err := news.HttpGet(thumbnailURL)
	if err != nil {
		return nil, err
	}

	if len(data) == 0 {
		return nil, nil
	}

	caption := ""
	if story["thumbnail"].(map[string]any)["caption"] != nil {
		caption = story["thumbnail"].(map[string]any)["caption"].(string)
	}

	return &news.Thumbnail{
		Image:   news.ConvertImage(data),
		Caption: news.SanitizeText(caption),
	}, nil
}
