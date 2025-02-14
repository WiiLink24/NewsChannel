package main

import "unicode/utf16"

// Headlines are the news articles that will appear on the News Channel banner in the Wii Menu.
type Headlines struct {
	HeadlineSize   uint32
	HeadlineOffset uint32
}

func (n *News) MakeWiiMenuHeadlines() {
	n.Header.HeadlinesTableOffset = n.GetCurrentSize()

	numberOfHeadlines := 11
	if len(n.articles) < 11 {
		numberOfHeadlines = len(n.articles)
	}

	n.Headlines = make([]Headlines, numberOfHeadlines)

	for i := 0; i < numberOfHeadlines; i++ {
		article := n.articles[i]

		// Encode to UTF-16
		encoded := utf16.Encode([]rune(article.Title))

		n.Headlines[i] = Headlines{
			HeadlineSize:   uint32(len(encoded)) * 2,
			HeadlineOffset: n.GetCurrentSize(),
		}

		n.HeadlineText = append(n.HeadlineText, encoded...)

		// Padding time.
		if (n.GetCurrentSize()+2)%4 == 0 {
			n.HeadlineText = append(n.HeadlineText, uint16(0))
		} else if (n.GetCurrentSize()+4)%4 == 0 {
			n.HeadlineText = append(n.HeadlineText, uint16(0))
			n.HeadlineText = append(n.HeadlineText, uint16(0))
		}
	}

	n.Header.NumberOfHeadlines = uint32(numberOfHeadlines)
}
