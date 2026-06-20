package reuters

import (
	"NewsChannel/news"
	"fmt"
)

func (r *Reuters) GetArticles() ([]news.Article, error) {
	var articles []news.Article

	if r.language == English {
		var categoryURLs = map[string]news.Topic{
			fmt.Sprintf("https://www.reuters.com/mobile/v1/world/%s/?outputType=json", r.country): news.NationalNews,
			"https://www.reuters.com/mobile/v1/world/?outputType=json":                            news.InternationalNews,
			"https://www.reuters.com/mobile/v1/sports/?outputType=json":                           news.Sports,
			"https://www.reuters.com/mobile/v1/lifestyle/?outputType=json":                        news.Entertainment,
			"https://www.reuters.com/mobile/v1/business/?outputType=json":                         news.Business,
			"https://www.reuters.com/mobile/v1/science/?outputType=json":                          news.Science,
			"https://www.reuters.com/mobile/v1/technology/?outputType=json":                       news.Technology,
		}
		for url, topic := range categoryURLs {
			temp, err := r.getMobileArticles(url, topic)
			if err != nil {
				return nil, err
			}

			articles = append(articles, temp...)
		}
	} else {
		var categoryURLs = map[string]news.Topic{
			"https://jp.reuters.com/pf/api/v3/content/fetch/articles-by-section-alias-or-id-v1?query={%22fetch_type%22:%22collection_or_section%22,%22orderby%22:%22last_updated_date:desc%22,%22section_id%22:%22/world/japan/%22,%22website%22:%22reuters-japan%22}":         news.NationalNews,
			"https://jp.reuters.com/pf/api/v3/content/fetch/articles-by-section-alias-or-id-v1?query={%22fetch_type%22:%22collection_or_section%22,%22orderby%22:%22last_updated_date:desc%22,%22section_id%22:%22/world/%22,%22website%22:%22reuters-japan%22}":               news.InternationalNews,
			"https://jp.reuters.com/pf/api/v3/content/fetch/articles-by-section-alias-or-id-v1?query={%22fetch_type%22:%22collection_or_section%22,%22orderby%22:%22last_updated_date:desc%22,%22section_id%22:%22/life/sports/%22,%22website%22:%22reuters-japan%22}":         news.Sports,
			"https://jp.reuters.com/pf/api/v3/content/fetch/articles-by-section-alias-or-id-v1?query={%22fetch_type%22:%22collection_or_section%22,%22orderby%22:%22last_updated_date:desc%22,%22section_id%22:%22/life/entertainment/%22,%22website%22:%22reuters-japan%22}":  news.Entertainment,
			"https://jp.reuters.com/pf/api/v3/content/fetch/articles-by-section-alias-or-id-v1?query={%22fetch_type%22:%22collection_or_section%22,%22orderby%22:%22last_updated_date:desc%22,%22section_id%22:%22/business/%22,%22website%22:%22reuters-japan%22}":            news.Business,
			"https://jp.reuters.com/pf/api/v3/content/fetch/articles-by-section-alias-or-id-v1?query={%22fetch_type%22:%22collection_or_section%22,%22orderby%22:%22last_updated_date:desc%22,%22section_id%22:%22/life/%22,%22website%22:%22reuters-japan%22}":                news.Science,
			"https://jp.reuters.com/pf/api/v3/content/fetch/articles-by-section-alias-or-id-v1?query={%22fetch_type%22:%22collection_or_section%22,%22orderby%22:%22last_updated_date:desc%22,%22section_id%22:%22/business/technology/%22,%22website%22:%22reuters-japan%22}": news.Technology,
		}
		for url, topic := range categoryURLs {
			temp, err := r.getWebArticles(url, topic)
			if err != nil {
				return nil, err
			}

			articles = append(articles, temp...)
		}
	}

	return articles, nil
}
