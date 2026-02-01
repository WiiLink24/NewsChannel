package ap

import (
	"NewsChannel/news"
	"fmt"
)

func (a *AP) GetArticles() ([]news.Article, error) {
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

	temp, err = a.GetScienceArticles()
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

func (a *AP) GetNationalArticles() ([]news.Article, error) {
	return a.getArticles(fmt.Sprintf("%s/apnews/topics/us-news", news.RSSHubAddress), news.NationalNews)
}

func (a *AP) GetInternationalArticles() ([]news.Article, error) {
	return a.getArticles(fmt.Sprintf("%s/apnews/topics/world-news", news.RSSHubAddress), news.InternationalNews)
}

func (a *AP) GetSportsArticles() ([]news.Article, error) {
	return a.getArticles(fmt.Sprintf("%s/apnews/topics/sports", news.RSSHubAddress), news.Sports)
}

func (a *AP) GetEntertainmentArticles() ([]news.Article, error) {
	return a.getArticles(fmt.Sprintf("%s/apnews/topics/entertainment", news.RSSHubAddress), news.Entertainment)
}

func (a *AP) GetBusinessArticles() ([]news.Article, error) {
	return a.getArticles(fmt.Sprintf("%s/apnews/topics/business", news.RSSHubAddress), news.Business)
}

func (a *AP) GetScienceArticles() ([]news.Article, error) {
	return a.getArticles(fmt.Sprintf("%s/apnews/topics/science", news.RSSHubAddress), news.Science)
}

func (a *AP) GetTechnologyArticles() ([]news.Article, error) {
	return a.getArticles(fmt.Sprintf("%s/apnews/topics/technology", news.RSSHubAddress), news.Technology)
}
