package main

import (
	"NewsChannel/news/reuters"
	"NewsChannel/news/rtve"
	_ "embed"
	"encoding/json"
	"fmt"
	"os"
	"time"
)

type Source struct {
	Logo            uint8
	Position        uint8
	_               uint16
	PictureSize     uint32
	PictureOffset   uint32
	NameSize        uint32
	NameOffset      uint32
	CopyrightSize   uint32
	CopyrightOffset uint32
}

func (n *News) GetNewsArticles() {
	// Choose news source based on country - RTVE for Spain, Reuters for others
	if n.currentCountry == "spain" {
		n.source = rtve.NewRTVE(n.oldArticleTitles)
	} else {
		n.source = reuters.NewReuters(n.oldArticleTitles, n.currentCountry)
	}

	var err error
	n.articles, err = n.source.GetArticles()
	if err != nil {
		panic(err)
	}

	// Save articles to file for inspection (Debug)
	// n.debugSaveArticles()
}

func (n *News) MakeSourceTable() {
	n.Header.SourceTableOffset = n.GetCurrentSize()

	logo := n.source.GetLogo()

	n.Sources = append(n.Sources, Source{
		Logo:            0,
		Position:        1,
		PictureSize:     uint32(len(logo)),
		PictureOffset:   0,
		NameSize:        0,
		NameOffset:      0,
		CopyrightSize:   0,
		CopyrightOffset: 0,
	})

	n.Sources[0].PictureOffset = n.GetCurrentSize()
	n.SourcePictures = logo

	for n.GetCurrentSize()%4 != 0 {
		n.SourcePictures = append(n.SourcePictures, 0)
	}

	n.Header.NumberOfSources = 1
}

// debugSaveArticles saves the fetched articles to a readable JSON file so you can see what was fetched.
func (n *News) debugSaveArticles() {
	if len(n.articles) == 0 {
		fmt.Printf("No articles found for country: %s\n", n.currentCountry)
		return
	}

	// Create directory
	err := os.MkdirAll("debug", 0755)
	if err != nil {
		fmt.Printf("Error creating debug directory: %v\n", err)
		return
	}

	// Structure 
	type DebugArticle struct {
		Title     string `json:"title"`
		Content   string `json:"content"`
		Topic     string `json:"topic"`
		Location  string `json:"location"`
		HasImage  bool   `json:"hasImage"`
		ImageSize int    `json:"imageSize"`
	}

	var debugArticles []DebugArticle
	topicNames := []string{"National", "International", "Sports", "Entertainment", "Business", "Science", "Technology"}

	for _, article := range n.articles {
		var content string
		if article.Content != nil {
			content = *article.Content
		} else {
			content = "No content"
		}

		var location string
		if article.Location != nil {
			location = article.Location.Name
		} else {
			location = "No location"
		}

		var topicName string
		if int(article.Topic) < len(topicNames) {
			topicName = topicNames[article.Topic]
		} else {
			topicName = fmt.Sprintf("Topic_%d", article.Topic)
		}

		var hasImage bool
		var imageSize int
		if article.Thumbnail != nil {
			hasImage = true
			imageSize = len(article.Thumbnail.Image)
		}

		debugArticles = append(debugArticles, DebugArticle{
			Title:     article.Title,
			Content:   content,
			Topic:     topicName,
			Location:  location,
			HasImage:  hasImage,
			ImageSize: imageSize,
		})
	}

	// Create filename with timestamp and country
	timestamp := time.Now().Format("2006-01-02_15-04-05")
	filename := fmt.Sprintf("debug/articles_%s_%s.json", n.currentCountry, timestamp)

	// Save to JSON file
	jsonData, err := json.MarshalIndent(debugArticles, "", "  ")
	if err != nil {
		fmt.Printf("Error marshaling articles: %v\n", err)
		return
	}

	err = os.WriteFile(filename, jsonData, 0644)
	if err != nil {
		fmt.Printf("Error writing debug file: %v\n", err)
		return
	}

	fmt.Printf("Debug: Saved %d articles to %s\n", len(debugArticles), filename)
}
