package ansa

import (
	"NewsChannel/news"
)

type Category struct {
	Name news.Topic
	URL  string
}

func (a *ANSA) GetArticles() ([]news.Article, error) {
	var categoryURLs = []Category{
		{
			news.NationalNews,
			"https://www.ansa.it/sito/ansait_rss.xml",
		},
		{
			news.InternationalNews,
			"https://www.ansa.it/sito/notizie/mondo/mondo_rss.xml",
		},
		{
			news.Sports,
			"https://www.ansa.it/sito/notizie/sport/sport_rss.xml",
		},
		{
			news.Entertainment,
			"https://www.ansa.it/sito/notizie/cultura/cultura_rss.xml",
		},
		{
			news.Business,
			"https://www.ansa.it/sito/notizie/economia/economia_rss.xml",
		},
		{
			news.Science,
			"https://www.ansa.it/canale_scienza_tecnica/notizie/scienzaetecnica_rss.xml",
		},
		{
			news.Technology,
			"https://www.ansa.it/canale_tecnologia/notizie/tecnologia_rss.xml",
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
