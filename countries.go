package main

func (n *News) GetTopicsForLanguage() []string {
	switch n.currentLanguageCode {
	case 0:
		return []string{"社会", "国際", "スポーツ", "芸能文化", "経済", "科学", "テクノロジー"}
	case 2:
		return []string{"Inland", "Weltweit", "Sport", "Feuilleton", "Wirtschaft", "Wissenschaft/Gesundheit", "Technologie"}
	case 3:
		return []string{"Actualités nationales", "Actualités internationales", "Sport", "Arts & loisirs", "Economie", "Sciences & santé", "Technologie"}
	case 4:
		return []string{"Nacional", "Internacional", "Deportes", "Cultura y ocio", "Economia", "Ciencia y salud", "Tecnologia"}
	case 5:
		return []string{"Notizie nazionali", "Notizie internazionali", "Sport", "Arte e Intrattenimento", "Attività commerciale", "Scienze e salute", "Tecnologia"}
	case 6:
		return []string{"Continentaal nieuws", "Internationaal nieuws", "Sport", "Kunst en entertainment", "Zakelljk nieuws", "Wetenschap en gezondheid", "Technologie"}
	default:
		return []string{"National News", "International News", "Sports", "Arts/Entertainment", "Business", "Science/Health", "Technology"}
	}
}
