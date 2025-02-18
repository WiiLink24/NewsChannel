package main

import (
	"NewsChannel/news/reuters"
	_ "embed"
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
	n.source = reuters.NewReuters(n.oldArticleTitles, n.currentCountry)

	var err error
	n.articles, err = n.source.GetArticles()
	if err != nil {
		panic(err)
	}
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
