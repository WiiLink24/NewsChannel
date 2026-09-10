package ap

import (
	"NewsChannel/news"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

var BaseUrl = "https://apnews.com/graphql/delivery/ap/v1"

func (a *AP) getArticles(path string, topic news.Topic) ([]news.Article, error) {
	parsedURL, err := url.Parse(BaseUrl)
	if err != nil {
		return nil, err
	}
	query := parsedURL.Query()
	query.Add("operationName", "ContentPageQuery")
	query.Add("variables", fmt.Sprintf("{\"path\":\"%s\"}", path))
	query.Add("extensions", "{\"persistedQuery\":{\"version\":1,\"sha256Hash\":\"3bc305abbf62e9e632403a74cc86dc1cba51156d2313f09b3779efec51fc3acb\"}}")
	parsedURL.RawQuery = query.Encode()
	requestURL := parsedURL.String()

	data, err := news.HttpGet(requestURL)
	if err != nil {
		return nil, err
	}

	var categoryListing ContentPageQuery
	err = json.Unmarshal(data, &categoryListing)
	if err != nil {
		return nil, err
	}

	var articles []news.Article
	for _, item := range categoryListing.Data.Screen.Main {
		if item.TypeName != "ColumnContainer" {
			continue
		}

		for _, column := range item.Columns {
			if column.TypeName != "PageListModule" {
				continue
			}

			for _, story := range column.Items {
				if story.TypeName != "PagePromo" {
					continue
				}

				title := news.SanitizeText(story.Title)
				// Check for duplicates
				if news.IsDuplicateArticle(a.oldArticleTitles, title) {
					continue
				}
				a.oldArticleTitles = append(a.oldArticleTitles, title)

				content, location, thumbnail, err := a.getFullArticle(story.Path)
				if err != nil {
					return nil, err
				}

				if content == "" {
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
				return articles, nil
			}
		}
	}

	return articles, nil
}

func (a *AP) getFullArticle(path string) (string, *news.Location, *news.Thumbnail, error) {
	if path == "" {
		return "", nil, nil, errors.New("empty path")
	}

	parsedURL, err := url.Parse(BaseUrl)
	if err != nil {
		return "", nil, nil, err
	}
	query := parsedURL.Query()
	query.Add("operationName", "StoryQuery")
	query.Add("variables", fmt.Sprintf("{\"path\":\"%s\"}", path))
	query.Add("extensions", "{\"persistedQuery\":{\"version\":1,\"sha256Hash\":\"d61508dd5f8c1d84fa6d49338e321266f56b4538c07c233564bb11300287d69c\"}}")
	parsedURL.RawQuery = query.Encode()
	requestURL := parsedURL.String()

	data, err := news.HttpGet(requestURL)
	if err != nil {
		return "", nil, nil, err
	}

	var articleResponse StoryPageQuery
	err = json.Unmarshal(data, &articleResponse)
	if err != nil {
		return "", nil, nil, err
	}

	// First, let's reconstruct the HTML that has been stored in the most horrific
	// format I've ever witnessed.
	var html strings.Builder
	for _, part := range articleResponse.Data.StoryPage.StoryBody {
		switch part.TypeName {
		case "HtmlString":
			html.WriteString(part.HTML)
		case "LinkEnhancement":
			for _, bodyPart := range part.Body {
				html.WriteString(bodyPart)
			}
		default:
			// We do not want other parts (video players, ads, etc)
			continue
		}
	}

	// Convert the now usable HTML to a string for use in the channel, and
	// extract the location from the first paragraph.
	content, locationString, err := a.extractArticleBody(html.String())
	if err != nil {
		return "", nil, nil, err
	}
	if len(content) == 0 {
		return "", nil, nil, nil
	}

	var location *news.Location
	if locationString != nil {
		location = news.GetLocationForExtractedLocation([]string{*locationString}, "en")
	}

	thumbnail := a.downloadThumbnail(articleResponse)

	return strings.TrimSpace(content), location, thumbnail, nil
}

func (a *AP) extractArticleBody(html string) (string, *string, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return "", nil, err
	}

	var contentSlice []string

	// Select the main article body div
	doc.Find("p").EachWithBreak(func(i int, elem *goquery.Selection) bool {
		text := strings.TrimSpace(elem.Text())
		if text != "" {
			contentSlice = append(contentSlice, text)
		}
		return true
	})

	if len(contentSlice) == 0 {
		return "", nil, nil
	}

	var content string
	for _, paragraph := range contentSlice {
		content += news.SanitizeText(paragraph)
		content += "\n\n"
	}

	// Get the location
	locationRegex := regexp.MustCompile(`(.*?) \(AP\) — `)
	location := locationRegex.FindStringSubmatch(contentSlice[0])
	if len(location) > 1 && len(location[1]) > 0 {
		return content, &location[1], nil
	}

	return content, nil, nil
}

func (a *AP) downloadThumbnail(articleResponse StoryPageQuery) *news.Thumbnail {
	var imageURL string
	var caption string

	if len(articleResponse.Data.StoryPage.BlendedGallery) > 0 {
		caption, imageURL = a.extractGalleryImage(articleResponse.Data.StoryPage.BlendedGallery)
	} else {
		caption, imageURL = a.extractLeadImage(articleResponse.Data.StoryPage.StoryLead)
	}

	var imageData []byte
	if imageURL == "" {
		return nil
	}

	imageData, err := news.HttpGet(imageURL)
	if err != nil || len(imageData) == 0 {
		return nil
	}

	return &news.Thumbnail{
		Image:   news.ConvertImage(imageData),
		Caption: news.SanitizeText(caption),
	}
}

func (a *AP) extractGalleryImage(gallery []BlendedGalleryElement) (caption string, imageURL string) {
	for _, item := range gallery {
		if item.TypeName != "Carousel" {
			continue
		}

		for _, slide := range item.Slides {
			if slide.TypeName != "GallerySlide" {
				continue
			}

			for _, captionString := range slide.Caption {
				if len(captionString) > 0 {
					caption = captionString
					break
				}
			}

			for _, media := range slide.Media {
				if media.TypeName != "Image" {
					continue
				}
				imageURL = a.extractURLFromMap(media.Image)
				return
			}
		}
	}
	return
}

func (a *AP) extractLeadImage(storyLead []StoryLeadElement) (caption string, imageURL string) {
	for _, item := range storyLead {
		if item.TypeName != "Figure" {
			continue
		}

		caption = item.AltText
		imageURL = a.extractURLFromMap(item.Image)
		return
	}

	return
}

func (a *AP) extractURLFromMap(imageMap ImageMap) string {
	if imageMap.TypeName != "Map" {
		return ""
	}
	for _, mapEntry := range imageMap.Entries {
		if mapEntry.TypeName != "MapEntry" || mapEntry.Key != "src" {
			continue
		}
		return mapEntry.Value
	}
	return ""
}
