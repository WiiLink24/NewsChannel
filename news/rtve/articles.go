package rtve

import (
	"NewsChannel/news"
	"log"
)

func (r *RTVE) GetArticles() ([]news.Article, error) {
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

func (r *RTVE) GetNationalArticles() ([]news.Article, error) {
	log.Printf("Fetching national articles from RTVE")
	url := "https://api.rtve.es/api/tematicas/1420/noticias.json?order=publication_date,desc"
	return r.getArticles(url, news.NationalNews)
}

func (r *RTVE) GetInternationalArticles() ([]news.Article, error) {
	url := "http://api.rtve.es/api/tematicas/828/noticias.json?order=publication_date,desc"
	return r.getArticles(url, news.InternationalNews)
}

func (r *RTVE) GetSportsArticles() ([]news.Article, error) {
	url := "http://www.rtve.es/servicios/scraper/?feed=/deportes&source=editorial&order=section"
	return r.getArticles(url, news.Sports)
}

func (r *RTVE) GetEntertainmentArticles() ([]news.Article, error) {
	url := "http://api.rtve.es/api/tematicas/827/noticias.json?order=publication_date,desc"
	return r.getArticles(url, news.Entertainment)
}

func (r *RTVE) GetBusinessArticles() ([]news.Article, error) {
	url := "http://api.rtve.es/api/tematicas/1011/noticias.json?order=publication_date,desc"
	return r.getArticles(url, news.Business)
}

func (r *RTVE) GetScienceArticles() ([]news.Article, error) {
	url := "http://api.rtve.es/api/tematicas/1012/noticias.json?order=publication_date,desc"
	return r.getArticles(url, news.Science)
}

func (r *RTVE) GetTechnologyArticles() ([]news.Article, error) {
	url := "http://api.rtve.es/api/tematicas/1161/noticias.json?order=publication_date,desc"
	return r.getArticles(url, news.Technology)
}
