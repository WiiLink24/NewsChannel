package tagesschau

import (
	"NewsChannel/news"
)

type Category struct {
	Name     news.Topic
	URL      string
	StoryKey string
}

func (r *Tagesschau) GetArticles() ([]news.Article, error) {
	var categoryURLs = []Category{
		{
			news.NationalNews,
			"https://www.tagesschau.de/api2u/news?ressort=inland",
			"news",
		},
		{
			news.InternationalNews,
			"https://www.tagesschau.de/api2u/news?ressort=ausland",
			"news",
		},
		{
			news.Sports,
			"https://www.tagesschau.de/api2u/news?ressort=sport",
			"news",
		},
		// For some reason, despite being a category on the site, entertainment does not exist as a category in the API
		// Search results work fine, it's just odd that you have to do it like this
		{
			news.Entertainment,
			"https://www.tagesschau.de/api2u/search?searchText=kultur",
			"searchResults",
		},
		{
			news.Business,
			"https://www.tagesschau.de/api2u/news?ressort=wirtschaft",
			"news",
		},
		{
			news.Science,
			"https://www.tagesschau.de/api2u/news?ressort=wissen",
			"news",
		},
		// Same deal as entertainment
		{
			news.Technology,
			"https://www.tagesschau.de/api2u/search?searchText=technologie",
			"searchResults",
		},
	}
	var articles []news.Article

	// Limited to 7 articles per file, more tend to push us over the file size limit.
	for len(articles) < 7 && len(categoryURLs) > 0 {
		var activeCategories []Category

		for _, cat := range categoryURLs {
			if len(articles) == 7 {
				break
			}

			var temp []news.Article
			var err error

			temp, err = r.getArticles(cat.URL, cat.Name, cat.StoryKey)

			if err != nil {
				return nil, err
			}

			if len(temp) > 0 {
				articles = append(articles, temp...)
				activeCategories = append(activeCategories, cat)
			}
		}

		categoryURLs = activeCategories
	}

	return articles, nil
}
