package tagesschau

import (
	"NewsChannel/news"
)

func (r *Tagesschau) GetArticles() ([]news.Article, error) {
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

func (r *Tagesschau) GetNationalArticles() ([]news.Article, error) {
	url := "https://www.tagesschau.de/api2u/news?ressort=inland"
	return r.getArticles(url, news.NationalNews, "news")
}

func (r *Tagesschau) GetInternationalArticles() ([]news.Article, error) {
	url := "https://www.tagesschau.de/api2u/news?ressort=ausland"
	return r.getArticles(url, news.InternationalNews, "news")
}

func (r *Tagesschau) GetSportsArticles() ([]news.Article, error) {
	url := "https://www.tagesschau.de/api2u/news?ressort=sport"
	return r.getArticles(url, news.Sports, "news")
}

func (r *Tagesschau) GetEntertainmentArticles() ([]news.Article, error) {
	// For some reason, despite being a category on the site, entertainment does not exist as a category in the API
	// Search results work fine, it's just odd that you have to do it like this
	url := "https://www.tagesschau.de/api2u/search?searchText=kultur"
	return r.getArticles(url, news.Entertainment, "searchResults")
}

func (r *Tagesschau) GetBusinessArticles() ([]news.Article, error) {
	url := "https://www.tagesschau.de/api2u/news?ressort=wirtschaft"
	return r.getArticles(url, news.Business, "news")
}

func (r *Tagesschau) GetScienceArticles() ([]news.Article, error) {
	url := "https://www.tagesschau.de/api2u/news?ressort=wissen"
	return r.getArticles(url, news.Science, "news")
}

func (r *Tagesschau) GetTechnologyArticles() ([]news.Article, error) {
	// Same deal as entertainment
	url := "https://www.tagesschau.de/api2u/search?searchText=technologie"
	return r.getArticles(url, news.Technology, "searchResults")
}
