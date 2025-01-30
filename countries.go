package main

func (n *News) GetTopicsForCountry() []string {
	switch n.currentCountryCode {
	case 18:
		return []string{"National News", "International News", "Sports", "Entertainment", "Business", "Science", "Technology"}
	}

	return []string{}
}
