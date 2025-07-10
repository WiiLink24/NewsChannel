package france24

import (
	"NewsChannel/news"
)

func (a *france24) GetArticles() ([]news.Article, error) {
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

func (a *france24) GetNationalArticles() ([]news.Article, error) {
	return a.getArticles("https://www.france24.com/fr/france/rss", news.NationalNews)
}

func (a *france24) GetInternationalArticles() ([]news.Article, error) {
	return a.getArticles("https://www.france24.com/fr/monde/rss", news.InternationalNews)
}

func (a *france24) GetSportsArticles() ([]news.Article, error) {
	return a.getArticles("https://www.france24.com/fr/sports/rss", news.Sports)
}

func (a *france24) GetEntertainmentArticles() ([]news.Article, error) {
	return a.getArticles("https://www.france24.com/fr/culture/rss", news.Entertainment)
}

func (a *france24) GetBusinessArticles() ([]news.Article, error) {
	return a.getArticles("https://www.france24.com/fr/economie/rss", news.Business)
}

func (a *france24) GetScienceArticles() ([]news.Article, error) {
	return a.getArticles("https://www.france24.com/fr/éco-tech/rss", news.Science)
}

func (a *france24) GetTechnologyArticles() ([]news.Article, error) {
	return a.getArticles("https://www.france24.com/fr/éco-tech/rss", news.Technology)
}
