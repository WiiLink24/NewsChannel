package reuters

import (
	"NewsChannel/news"
	"fmt"
)

func (r *Reuters) GetArticles() ([]news.Article, error) {
	var articles []news.Article

	temp, err := r.GetNationalArticles()
	if err != nil {
		return nil, err
	}

	articles = append(articles, temp...)

	temp, err = r.GetInternationalArticles()
	if err != nil {
		return nil, err
	}

	articles = append(articles, temp...)

	temp, err = r.GetSportsArticles()
	if err != nil {
		return nil, err
	}

	articles = append(articles, temp...)

	temp, err = r.GetEntertainmentArticles()
	if err != nil {
		return nil, err
	}

	articles = append(articles, temp...)

	temp, err = r.GetBusinessArticles()
	if err != nil {
		return nil, err
	}

	articles = append(articles, temp...)

	temp, err = r.GetScienceArticles()
	if err != nil {
		return nil, err
	}

	articles = append(articles, temp...)

	temp, err = r.GetTechnologyArticles()
	if err != nil {
		return nil, err
	}

	articles = append(articles, temp...)

	return articles, nil
}

func (r *Reuters) GetNationalArticles() ([]news.Article, error) {
	url := fmt.Sprintf("https://www.reuters.com/arc/outboundfeeds/v4/mobile/section/world/%s/?outputType=json&from=0&size=20", r.country)
	return getArticles(url, news.NationalNews)
}

func (r *Reuters) GetInternationalArticles() ([]news.Article, error) {
	url := "https://www.reuters.com/arc/outboundfeeds/v4/mobile/section/world/?outputType=json&from=0&size=20"
	return getArticles(url, news.InternationalNews)
}

func (r *Reuters) GetSportsArticles() ([]news.Article, error) {
	url := "https://www.reuters.com/arc/outboundfeeds/v4/mobile/section/sports/?outputType=json&from=0&size=20"
	return getArticles(url, news.Sports)
}

func (r *Reuters) GetEntertainmentArticles() ([]news.Article, error) {
	url := "https://www.reuters.com/arc/outboundfeeds/v4/mobile/section/lifestyle/?outputType=json&from=0&size=20"
	return getArticles(url, news.Entertainment)
}

func (r *Reuters) GetBusinessArticles() ([]news.Article, error) {
	url := "https://www.reuters.com/arc/outboundfeeds/v4/mobile/section/business/?outputType=json&from=0&size=20"
	return getArticles(url, news.Business)
}

func (r *Reuters) GetScienceArticles() ([]news.Article, error) {
	url := "https://www.reuters.com/arc/outboundfeeds/v4/mobile/section/science/?outputType=json&from=0&size=20"
	return getArticles(url, news.Science)
}

func (r *Reuters) GetTechnologyArticles() ([]news.Article, error) {
	url := "https://www.reuters.com/arc/outboundfeeds/v4/mobile/section/technology/?outputType=json&from=0&size=20"
	return getArticles(url, news.Technology)
}
