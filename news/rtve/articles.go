package rtve

import (
	"NewsChannel/news"
)

type Category struct {
	Name news.Topic
	ID   string
}

func (r *RTVE) GetArticles() ([]news.Article, error) {
	var categoryURLs = []Category{
		{
			news.NationalNews,
			"1420",
		},
		{
			news.InternationalNews,
			"828",
		},
		{
			news.Sports,
			"816",
		},
		{
			news.Entertainment,
			"827",
		},
		{
			news.Business,
			"1011",
		},
		{
			news.Science,
			"1012",
		},
		{
			news.Technology,
			"1161",
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

			temp, err = r.getArticles(cat.ID, cat.Name)

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
