package nos

import (
	"NewsChannel/news"
)

type Category struct {
	Name news.Topic
	URL  string
}

func (a *nos) GetArticles() ([]news.Article, error) {
	var categoryURLs = []Category{
		{
			news.NationalNews,
			"nosnieuwsbinnenland",
		},
		{
			news.InternationalNews,
			"nosnieuwsbuitenland",
		},
		{
			news.Sports,
			"nossportalgemeen",
		},
		{
			news.Entertainment,
			"nosnieuwsopmerkelijk",
		},
		{
			news.Business,
			"nosnieuwseconomie",
		},
		{
			news.Science,
			"nosnieuwscultuurenmedia",
		},
		{
			news.Technology,
			"nosnieuwstech",
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

			temp, err = a.getArticles(cat.URL, cat.Name)

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
