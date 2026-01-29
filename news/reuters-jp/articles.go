package reutersjp

import (
	"NewsChannel/news"
)

func (r *ReutersJP) GetArticles() ([]news.Article, error) {
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

func (r *ReutersJP) GetNationalArticles() ([]news.Article, error) {
	url := "https://jp.reuters.com/pf/api/v3/content/fetch/articles-by-section-alias-or-id-v1?query={%22fetch_type%22:%22collection_or_section%22,%22orderby%22:%22last_updated_date:desc%22,%22section_id%22:%22/world/japan/%22,%22website%22:%22reuters-japan%22}"
	return r.getArticles(url, news.NationalNews)
}

func (r *ReutersJP) GetInternationalArticles() ([]news.Article, error) {
	url := "https://jp.reuters.com/pf/api/v3/content/fetch/articles-by-section-alias-or-id-v1?query={%22fetch_type%22:%22collection_or_section%22,%22orderby%22:%22last_updated_date:desc%22,%22section_id%22:%22/world/%22,%22website%22:%22reuters-japan%22}"
	return r.getArticles(url, news.InternationalNews)
}

func (r *ReutersJP) GetSportsArticles() ([]news.Article, error) {
	url := "https://jp.reuters.com/pf/api/v3/content/fetch/articles-by-section-alias-or-id-v1?query={%22fetch_type%22:%22collection_or_section%22,%22orderby%22:%22last_updated_date:desc%22,%22section_id%22:%22/life/sports/%22,%22website%22:%22reuters-japan%22}"
	return r.getArticles(url, news.Sports)
}

func (r *ReutersJP) GetEntertainmentArticles() ([]news.Article, error) {
	url := "https://jp.reuters.com/pf/api/v3/content/fetch/articles-by-section-alias-or-id-v1?query={%22fetch_type%22:%22collection_or_section%22,%22orderby%22:%22last_updated_date:desc%22,%22section_id%22:%22/life/entertainment/%22,%22website%22:%22reuters-japan%22}"
	return r.getArticles(url, news.Entertainment)
}

func (r *ReutersJP) GetBusinessArticles() ([]news.Article, error) {
	url := "https://jp.reuters.com/pf/api/v3/content/fetch/articles-by-section-alias-or-id-v1?query={%22fetch_type%22:%22collection_or_section%22,%22orderby%22:%22last_updated_date:desc%22,%22section_id%22:%22/business/%22,%22website%22:%22reuters-japan%22}"
	return r.getArticles(url, news.Business)
}

func (r *ReutersJP) GetScienceArticles() ([]news.Article, error) {
	url := "https://jp.reuters.com/pf/api/v3/content/fetch/articles-by-section-alias-or-id-v1?query={%22fetch_type%22:%22collection_or_section%22,%22orderby%22:%22last_updated_date:desc%22,%22section_id%22:%22/life/%22,%22website%22:%22reuters-japan%22}"
	return r.getArticles(url, news.Science)
}

func (r *ReutersJP) GetTechnologyArticles() ([]news.Article, error) {
	url := "https://jp.reuters.com/pf/api/v3/content/fetch/articles-by-section-alias-or-id-v1?query={%22fetch_type%22:%22collection_or_section%22,%22orderby%22:%22last_updated_date:desc%22,%22section_id%22:%22/business/technology/%22,%22website%22:%22reuters-japan%22}"
	return r.getArticles(url, news.Technology)
}
