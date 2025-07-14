package ansa

import (
	"NewsChannel/news"
)

func (a *ANSA) GetArticles() ([]news.Article, error) {
	var articles []news.Article

	temp, err := a.GetNationalArticles()
	if err != nil {
		return nil, err
	}
	articles = append(articles, temp...)

	temp, err = a.GetInternationalArticles()
	if err != nil {
		return nil, err
	}
	articles = append(articles, temp...)

	temp, err = a.GetSportsArticles()
	if err != nil {
		return nil, err
	}
	articles = append(articles, temp...)

	temp, err = a.GetEntertainmentArticles()
	if err != nil {
		return nil, err
	}
	articles = append(articles, temp...)

	temp, err = a.GetBusinessArticles()
	if err != nil {
		return nil, err
	}
	articles = append(articles, temp...)

	temp, err = a.GetTechnologyArticles()
	if err != nil {
		return nil, err
	}
	articles = append(articles, temp...)

	return articles, nil
}

func (a *ANSA) GetNationalArticles() ([]news.Article, error) {
	return a.getArticles("https://www.ansa.it/sito/ansait_rss.xml", news.NationalNews)
}

func (a *ANSA) GetInternationalArticles() ([]news.Article, error) {
	return a.getArticles("https://www.ansa.it/sito/notizie/mondo/mondo_rss.xml", news.InternationalNews)
}

func (a *ANSA) GetSportsArticles() ([]news.Article, error) {
	return a.getArticles("https://www.ansa.it/sito/notizie/sport/sport_rss.xml", news.Sports)
}

func (a *ANSA) GetEntertainmentArticles() ([]news.Article, error) {
	return a.getArticles("https://www.ansa.it/sito/notizie/cultura/cultura_rss.xml", news.Entertainment)
}

func (a *ANSA) GetBusinessArticles() ([]news.Article, error) {
	return a.getArticles("https://www.ansa.it/sito/notizie/economia/economia_rss.xml", news.Business)
}

func (a *ANSA) GetScienceArticles() ([]news.Article, error) {
	// ANSA doesn't have a separate science feed, use technology
	return a.getArticles("https://www.ansa.it/canale_tecnologia/notizie/tecnologia_rss.xml", news.Science)
}

func (a *ANSA) GetTechnologyArticles() ([]news.Article, error) {
	return a.getArticles("https://www.ansa.it/canale_tecnologia/notizie/tecnologia_rss.xml", news.Technology)
}
