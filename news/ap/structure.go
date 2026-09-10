package ap

type ContentPageQuery struct {
	Data struct {
		Screen struct {
			Main []struct {
				TypeName string            `json:"__typename"`
				Columns  []ColumnContainer `json:"columns"`
			} `json:"main"`
		} `json:"Screen"`
	} `json:"data"`
}

type ColumnContainer struct {
	TypeName string           `json:"__typename"`
	Items    []ArticleSummary `json:"items"`
}

type ArticleSummary struct {
	TypeName string `json:"__typename"`
	Title    string `json:"title"`
	Path     string `json:"graphqlPath"`
}

type StoryPageQuery struct {
	Data struct {
		StoryPage struct {
			TypeName       string                  `json:"__typename"`
			StoryLead      []StoryLeadElement      `json:"storyLead"`
			BlendedGallery []BlendedGalleryElement `json:"blendedGallery"`
			StoryBody      []StoryElement          `json:"storyBody"`
		} `json:"StoryPage"`
	} `json:"data"`
}

type StoryLeadElement struct {
	TypeName string   `json:"__typename"`
	AltText  string   `json:"alt"`
	Image    ImageMap `json:"image"`
}

type BlendedGalleryElement struct {
	TypeName string `json:"__typename"`
	Slides   []struct {
		TypeName string   `json:"__typename"`
		Caption  []string `json:"caption"`
		Media    []struct {
			TypeName string   `json:"__typename"`
			Image    ImageMap `json:"image"`
		} `json:"media"`
	} `json:"slides"`
}

type ImageMap struct {
	TypeName string `json:"__typename"`
	Entries  []struct {
		TypeName string `json:"__typename"`
		Key      string `json:"key"`
		Value    string `json:"value"`
	} `json:"entries"`
}

type StoryElement struct {
	TypeName string   `json:"__typename"`
	HTML     string   `json:"html"`
	Body     []string `json:"body"`
}
