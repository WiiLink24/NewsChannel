package main

type Header struct {
	Version          uint32
	Filesize         uint32
	CRC32            uint32
	UpdatedTimestamp uint32
	EndTimestamp     uint32
	// The country code is written as little endian and I really don't want to deal with that.
	// No country code should overflow an u8 so this is fine.
	CountryCode              uint8
	_                        [3]byte
	UpdatedTimestamp2        uint32
	SupportedLanguages       [16]uint8
	LanguageCode             uint8
	GooFlag                  uint8
	ShowLanguageSelectScreen uint8
	DownloadInterval         uint8
	MessageOffset            uint32
	NumberOfTopics           uint32
	TopicTableOffset         uint32
	NumberOfArticles         uint32
	ArticleTableOffset       uint32
	NumberOfSources          uint32
	SourceTableOffset        uint32
	NumberOfLocations        uint32
	LocationTableOffset      uint32
	NumberOfImages           uint32
	ImagesTableOffset        uint32
	DownloadCount            uint16
	_                        uint16
	NumberOfHeadlines        uint32
	HeadlinesTableOffset     uint32
}

func (n *News) MakeHeader() {
	n.Header = Header{
		Version:                  512,
		Filesize:                 0,
		CRC32:                    0,
		UpdatedTimestamp:         fixTime(currentTime),
		EndTimestamp:             fixTime(currentTime) + 1500,
		CountryCode:              n.currentCountryCode,
		UpdatedTimestamp2:        fixTime(currentTime),
		SupportedLanguages:       [16]uint8{1, 3, 4, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF},
		LanguageCode:             n.currentLanguageCode,
		GooFlag:                  0,
		ShowLanguageSelectScreen: 0,
		DownloadInterval:         30,
		MessageOffset:            0,
		NumberOfTopics:           0,
		TopicTableOffset:         0,
		NumberOfArticles:         0,
		ArticleTableOffset:       0,
		NumberOfSources:          0,
		SourceTableOffset:        0,
		NumberOfLocations:        0,
		LocationTableOffset:      0,
		NumberOfImages:           0,
		ImagesTableOffset:        0,
		DownloadCount:            480,
		NumberOfHeadlines:        0,
		HeadlinesTableOffset:     0,
	}
}
