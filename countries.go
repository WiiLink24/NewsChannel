package main

func (n *News) GetTopicsForLanguage() []string {
	switch n.currentLanguageCode {
	case 0:
		return []string{"社会", "国際", "スポーツ", "芸能文化", "経済", "科学", "テクノロジー"}
	case 2:
		return []string{"Deutschland", "Weltnachrichten", "Sport", "Unterhaltung", "Wirtschaft", "Gesundheit/Medizin", "Technik"}
	case 3:
		return []string{"France", "Monde", "Sports", "Culture", "Economie", "Sciences & santé", "Technologies"}
	case 4:
		return []string{"España", "Internacional", "Deportes", "Cultura", "Economia", "Ciencia y salud", "Tecnologia"}
	case 5:
		return []string{"Notizie nazionali", "Notizie internazionali", "Sport", "Arte e Intrattenimento", "Attività commerciale", "Scienze e salute", "Tecnologia"}
	case 6:
		return []string{"Continentaal nieuws", "Internationaal nieuws", "Sport", "Kunst en entertainment", "Zakelljk nieuws", "Wetenschap en gezondheid", "Technologie"}
	default:
		return []string{"National News", "International News", "Sports", "Entertainment", "Business", "Science", "Technology"}
	}
}
