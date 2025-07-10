package welt

import (
	"NewsChannel/news"
)

func (a *welt) GetArticles() ([]news.Article, error) {
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

func (a *welt) GetNationalArticles() ([]news.Article, error) {
	return a.getArticles("https://www.welt.de/feeds/topnews.rss", news.NationalNews)
}

func (a *welt) GetInternationalArticles() ([]news.Article, error) {
	return a.getArticles("https://www.welt.de/feeds/latest.rss", news.InternationalNews)
}

func (a *welt) GetSportsArticles() ([]news.Article, error) {
	return a.getArticles("https://www.welt.de/feeds/section/sport.rss", news.Sports)
}

func (a *welt) GetEntertainmentArticles() ([]news.Article, error) {
	return a.getArticles("https://www.welt.de/feeds/section/kultur.rss", news.Entertainment)
}

func (a *welt) GetBusinessArticles() ([]news.Article, error) {
	return a.getArticles("https://www.welt.de/feeds/section/wirtschaft.rss", news.Business)
}

func (a *welt) GetScienceArticles() ([]news.Article, error) {
	return a.getArticles("https://www.welt.de/feeds/section/wissenschaft.rss", news.Science)
}

func (a *welt) GetTechnologyArticles() ([]news.Article, error) {
	return a.getArticles("https://www.welt.de/feeds/section/motor.rss", news.Technology)
}
