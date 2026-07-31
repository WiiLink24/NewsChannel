package reuters

import (
	"NewsChannel/news"
	"fmt"
)

type Category struct {
	Name      news.Topic
	JapanURL  string
	MobileURL string
}

func (r *Reuters) GetArticles() ([]news.Article, error) {
	var categoryURLs = []Category{
		{
			news.NationalNews,
			"/world/japan/",
			fmt.Sprintf("/world/%s/", r.country),
		},
		{
			news.InternationalNews,
			"/world/",
			"/world/",
		},
		{
			news.Sports,
			"/life/sports/",
			"/sports/",
		},
		{
			news.Entertainment,
			"/life/entertainment/",
			"/lifestyle/",
		},
		{
			news.Business,
			"/business/",
			"/business/",
		},
		{
			news.Science,
			"/life/",
			"/science/",
		},
		{
			news.Technology,
			"/business/technology/",
			"/technology/",
		},
	}
	var articles []news.Article

	// Limited to 7 articles per file, more tend to push us over the file size limit.
	for len(articles) < 7 && len(categoryURLs) > 0 {
		var activeCategories []Category

		for _, cat := range categoryURLs {
			if len(articles) == 7 {
				break
			}

			var temp []news.Article
			var err error

			if r.language == English {
				temp, err = r.getMobileArticles(cat.MobileURL, cat.Name)
			} else {
				temp, err = r.getWebArticles(cat.JapanURL, cat.Name)
			}

			if err != nil {
				return nil, err
			}

			if len(temp) > 0 {
				articles = append(articles, temp...)
				activeCategories = append(activeCategories, cat)
			}
		}

		categoryURLs = activeCategories
	}

	return articles, nil
}
