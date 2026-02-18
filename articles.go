package main

import (
	"math"
	"unicode/utf16"
)

type Article struct {
	ID                uint32
	SourceIndex       uint32
	LocationIndex     uint32
	PictureTimestamp  uint32
	PictureIndex      uint32
	PublishedTime     uint32
	UpdatedTime       uint32
	HeadlineSize      uint32
	HeadlineOffset    uint32
	ArticleTextSize   uint32
	ArticleTextOffset uint32
}

type Image struct {
	CreditSize    uint32
	CreditOffset  uint32
	CaptionSize   uint32
	CaptionOffset uint32
	PictureSize   uint32
	PictureOffset uint32
}

func (n *News) MakeArticleTable() {
	n.Header.ArticleTableOffset = n.GetCurrentSize()

	// First write all metadata
	for i, article := range n.articles {
		publishedTime := currentTime

		// Parse the location if any.
		locationIndex := uint32(math.MaxUint32)
		if article.Location != nil {
			for i, location := range n.locations {
				if floatCompare(location.Latitude, article.Location.Latitude) && floatCompare(location.Longitude, article.Location.Longitude) {
					locationIndex = uint32(i)
					break
				}
			}

			// If no existing location was found but the article contains a location, determine if we need to make an API call.
			if locationIndex == uint32(math.MaxUint32) && article.Location.Latitude != 0 {
				locationIndex = uint32(len(n.locations))
				n.locations = append(n.locations, article.Location)
			} else if locationIndex == uint32(math.MaxUint32) && article.Location.Name != "" {
				// TODO: API Call for location
			}
		}

		n.Articles = append(n.Articles, Article{
			ID:                uint32(i + 1),
			SourceIndex:       0,
			LocationIndex:     locationIndex,
			PictureTimestamp:  0,
			PictureIndex:      math.MaxUint32,
			PublishedTime:     fixTime(publishedTime),
			UpdatedTime:       fixTime(currentTime),
			HeadlineSize:      0,
			HeadlineOffset:    0,
			ArticleTextSize:   0,
			ArticleTextOffset: 0,
		})

		n.timestamps[article.Topic+1] = append(n.timestamps[article.Topic+1], Timestamp{
			Time:          fixTime(currentTime),
			ArticleNumber: uint32(i + 1),
		})
	}

	// Next write the text
	for i, article := range n.articles {
		encodedTitle := utf16.Encode([]rune(article.Title))
		encodedArticle := utf16.Encode([]rune(*article.Content))

		n.Articles[i].HeadlineSize = uint32(len(encodedTitle) * 2)
		n.Articles[i].ArticleTextSize = uint32(len(encodedArticle) * 2)

		n.Articles[i].HeadlineOffset = n.GetCurrentSize()
		n.ArticleText = append(n.ArticleText, encodedTitle...)

		// Null terminator
		n.ArticleText = append(n.ArticleText, 0)

		for n.GetCurrentSize()%4 != 0 {
			n.ArticleText = append(n.ArticleText, uint16(0))
		}

		n.Articles[i].ArticleTextOffset = n.GetCurrentSize()
		n.ArticleText = append(n.ArticleText, encodedArticle...)

		// Null terminator
		n.ArticleText = append(n.ArticleText, 0)

		for n.GetCurrentSize()%4 != 0 {
			n.ArticleText = append(n.ArticleText, uint16(0))
		}
	}

	n.Header.NumberOfArticles = uint32(len(n.Articles))
}

func (n *News) WriteImages() {
	n.Header.ImagesTableOffset = n.GetCurrentSize()
	for _, article := range n.articles {
		if article.Thumbnail == nil || len(article.Thumbnail.Image) == 0 {
			continue
		}

		n.Images = append(n.Images, Image{
			CreditSize:    0,
			CreditOffset:  0,
			CaptionSize:   0,
			CaptionOffset: 0,
			PictureSize:   uint32(len(article.Thumbnail.Image)),
			PictureOffset: 0,
		})
	}

	i := 0
	for j, article := range n.articles {
		if article.Thumbnail == nil || len(article.Thumbnail.Image) == 0 {
			continue
		}

		n.Images[i].PictureOffset = n.GetCurrentSize()
		n.ImagesData = append(n.ImagesData, article.Thumbnail.Image...)
		for n.GetCurrentSize()%4 != 0 {
			n.ImagesData = append(n.ImagesData, 0)
		}

		// Fix up the article
		n.Articles[j].PictureIndex = uint32(i)
		n.Articles[j].PictureTimestamp = fixTime(currentTime)
		i++
	}

	i = 0
	for _, article := range n.articles {
		if article.Thumbnail == nil {
			continue
		}
		if len(article.Thumbnail.Caption) == 0 {
			i++
			continue
		}

		caption := utf16.Encode([]rune(article.Thumbnail.Caption))
		n.Images[i].CaptionOffset = n.GetCurrentSize()
		n.Images[i].CaptionSize = uint32(len(caption) / 2)
		n.CaptionData = append(n.CaptionData, caption...)
		n.CaptionData = append(n.CaptionData, 0)

		for n.GetCurrentSize()%4 != 0 {
			n.CaptionData = append(n.CaptionData, uint16(0))
		}

		i++
	}

	n.Header.NumberOfImages = uint32(len(n.Images))
}
