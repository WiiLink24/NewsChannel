package reuters

import (
	"NewsChannel/news"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
)

var htmlRegex = regexp.MustCompile("<.*?>")

func httpGet(url string) ([]byte, error) {
	c := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	// User-Agent derived from the iPadOS version of the Reuters App. Bypasses the no JS screen.
	req.Header.Set("User-Agent", "ReutersNews/7.6.0 iPad8,6 iPadOS/18.1 CFNetwork/1.0 Darwin/24.1.0")

	resp, err := c.Do(req)
	if err != nil || resp.StatusCode != http.StatusOK {
		return nil, err
	}

	// Read the body
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func (r *Reuters) getArticles(url string, topic news.Topic) ([]news.Article, error) {
	data, err := httpGet(url)
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
		if v["type"].(string) != "latest-stories" {
			continue
		}

		stories := v["data"].(map[string]any)["stories"].([]any)
		i := 0
		for _, story := range stories {
			if i == 1 {
				break
			}

			title := story.(map[string]any)["title"].(string)
			// Compare previous articles to see if we have a duplicate.
			if news.IsDuplicateArticle(r.oldArticleTitles, title) {
				continue
			}

			// The article is nested inside a "templates" list, with the data we require in the 1st index.
			// I (Noah) refer to this as bad because it returns the web page, rather than the mobile API page.
			// The mobile API is much easier to parse.
			articlePath := story.(map[string]any)["url"]
			articleURL := fmt.Sprintf("https://www.reuters.com/mobile/v1%s", articlePath)
			articleData, err := httpGet(articleURL)
			if err != nil {
				return nil, err
			}

			// Parse article JSON
			var articleJSON []map[string]any
			err = json.Unmarshal(articleData, &articleJSON)
			if err != nil {
				if err.Error() == "invalid character '<' looking for beginning of value" {
					continue
				}

				return nil, err
			}

			content, err := parseArticle(articleJSON)
			if err != nil {
				return nil, err
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
			i++
		}
	}

	return articles, nil
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

			// Remove any <a> strings.
			unSanitized := content.(map[string]any)["content"].(string)
			sanitized := htmlRegex.ReplaceAllString(unSanitized, "")
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

		thumbnailURL := child["data"].(map[string]any)["article"].(map[string]any)["thumbnail"].(map[string]any)["resizer_url"].(string)

		// Add the required params
		thumbnailURL = fmt.Sprintf("%s&width=200&height=200", thumbnailURL)
		data, err := httpGet(thumbnailURL)
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

		if child["data"].(map[string]any)["article"].(map[string]any)["dateline"] == nil {
			return nil, nil
		}

		location := child["data"].(map[string]any)["article"].(map[string]any)["dateline"].([]any)[0].(string)
		location = strings.ToUpper(strings.Split(location, ",")[0])

		if l, ok := news.CommonLocations[location]; ok {
			return &l, nil
		}

		// Return just the name as we don't want to make API calls yet.
		return &news.Location{
			Name: location,
		}, nil
	}

	return nil, nil
}
