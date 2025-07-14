package nos

import (
	"NewsChannel/news"
)

func (a *nos) GetArticles() ([]news.Article, error) {
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

func (a *nos) GetNationalArticles() ([]news.Article, error) {
	return a.getArticles("https://feeds.nos.nl/nosnieuwsbinnenland", news.NationalNews)
}

func (a *nos) GetInternationalArticles() ([]news.Article, error) {
	return a.getArticles("https://feeds.nos.nl/nosnieuwsbuitenland", news.InternationalNews)
}

func (a *nos) GetSportsArticles() ([]news.Article, error) {
	return a.getArticles("https://feeds.nos.nl/nossportalgemeen", news.Sports)
}

func (a *nos) GetEntertainmentArticles() ([]news.Article, error) {
	return a.getArticles("https://feeds.nos.nl/nosnieuwsopmerkelijk", news.Entertainment)
}

func (a *nos) GetBusinessArticles() ([]news.Article, error) {
	return a.getArticles("https://feeds.nos.nl/nosnieuwseconomie", news.Business)
}

func (a *nos) GetScienceArticles() ([]news.Article, error) {
	return a.getArticles("https://feeds.nos.nl/nosnieuwscultuurenmedia", news.Science)
}

func (a *nos) GetTechnologyArticles() ([]news.Article, error) {
	return a.getArticles("https://feeds.nos.nl/nosnieuwstech", news.Technology)
}
