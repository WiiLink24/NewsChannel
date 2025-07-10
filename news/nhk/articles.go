package nhk

import (
	"NewsChannel/news"
)

func (a *nhk) GetArticles() ([]news.Article, error) {
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

func (a *nhk) GetNationalArticles() ([]news.Article, error) {
	return a.getArticles("https://www.nhk.or.jp/rss/news/cat0.xml", news.NationalNews)
}

func (a *nhk) GetInternationalArticles() ([]news.Article, error) {
	return a.getArticles("https://www.nhk.or.jp/rss/news/cat6.xml", news.InternationalNews)
}

func (a *nhk) GetSportsArticles() ([]news.Article, error) {
	return a.getArticles("https://www.nhk.or.jp/rss/news/cat7.xml", news.Sports)
}

func (a *nhk) GetEntertainmentArticles() ([]news.Article, error) {
	return a.getArticles("https://www.nhk.or.jp/rss/news/cat2.xml", news.Entertainment)
}

func (a *nhk) GetBusinessArticles() ([]news.Article, error) {
	return a.getArticles("https://www.nhk.or.jp/rss/news/cat5.xml", news.Business)
}

func (a *nhk) GetScienceArticles() ([]news.Article, error) {
	return a.getArticles("https://www.nhk.or.jp/rss/news/cat3.xml", news.Science)
}

func (a *nhk) GetTechnologyArticles() ([]news.Article, error) {
	// This is actually "Society", please sketch change the name of the topic if possible :3
	return a.getArticles("https://www.nhk.or.jp/rss/news/cat1.xml", news.Technology)
}
