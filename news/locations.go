package news

type Location struct {
	Longitude float64
	Latitude  float64
	Name      string
}

var CommonLocations = map[string]Location{
	"AMSTERDAM": {
		Longitude: 4.883423,
		Latitude:  52.366333,
		Name:      "Amsterdam",
	},
	"ATLANTA": {
		Longitude: -84.385986,
		Latitude:  33.744507,
		Name:      "Atlanta",
	},
	"BAGHDAD": {
		Longitude: 44.412231,
		Latitude:  33.348999,
		Name:      "Baghdad",
	},
	"BALTIMORE": {
		Longitude: -76.607666,
		Latitude:  39.287109,
		Name:      "Baltimore",
	},
	"BANGKOK": {
		Longitude: 100.513916,
		Latitude:  13.749390,
		Name:      "Bangkok",
	},
	"BEIJING": {
		Longitude: 116.433105,
		Latitude:  39.913330,
		Name:      "Beijing",
	},
	"BEIRUT": {
		Longitude: 35.496826,
		Latitude:  33.881836,
		Name:      "Beirut",
	},
	"BERLIN": {
		Longitude: 13.403320,
		Latitude:  52.520142,
		Name:      "Berlin",
	},
	"BOSTON": {
		Longitude: -71.059570,
		Latitude:  42.357788,
		Name:      "Boston",
	},
	"BRUSSELS": {
		Longitude: 4.367065,
		Latitude:  50.839233,
		Name:      "Brussels",
	},
	"CAIRO": {
		Longitude: 31.245117,
		Latitude:  30.047607,
		Name:      "Cairo",
	},
	"CHICAGO": {
		Longitude: -87.648926,
		Latitude:  41.846924,
		Name:      "Chicago",
	},
	"CINCINNATI": {
		Longitude: -84.451904,
		Latitude:  39.160767,
		Name:      "Cincinnati",
	},
	"CLEVELAND": {
		Longitude: -81.694336,
		Latitude:  41.495361,
		Name:      "Cleveland",
	},
	"DALLAS": {
		Longitude: -96.795044,
		Latitude:  32.783203,
		Name:      "Dallas",
	},
	"DENVER": {
		Longitude: -104.979858,
		Latitude:  39.737549,
		Name:      "Denver",
	},
	"DETROIT": {
		Longitude: -83.045654,
		Latitude:  42.330322,
		Name:      "Detroit",
	},
	"DJIBOUTI": {
		Longitude: 43.148804,
		Latitude:  11.596069,
		Name:      "Djibouti",
	},
	"DUBLIN": {
		Longitude: -6.225898,
		Latitude:  53.366550,
		Name:      "Dublin",
	},
	"GENEVA": {
		Longitude: 6.168823,
		Latitude:  46.197510,
		Name:      "Geneva",
	},
	"GIBRALTAR": {
		Longitude: -5.345272,
		Latitude:  36.121167,
		Name:      "Gibraltar",
	},
	"GUATEMALA CITY": {
		Longitude: -90.521851,
		Latitude:  14.617310,
		Name:      "Guatemala City",
	},
	"HAVANA": {
		Longitude: -82.348022,
		Latitude:  23.148193,
		Name:      "Havana",
	},
	"HELSINKI": {
		Longitude: 24.933472,
		Latitude:  60.166626,
		Name:      "Helsinki",
	},
	"HONG KONG": {
		Longitude: 114.296265,
		Latitude:  22.461548,
		Name:      "Hong Kong",
	},
	"HONOLULU": {
		Longitude: -157.857056,
		Latitude:  21.302490,
		Name:      "Honolulu",
	},
	"HOUSTON": {
		Longitude: -95.361328,
		Latitude:  29.761963,
		Name:      "Houston",
	},
	"INDIANAPOLIS": {
		Longitude: -86.154785,
		Latitude:  39.765015,
		Name:      "Indianapolis",
	},
	"ISLAMABAD": {
		Longitude: 73.163452,
		Latitude:  33.695068,
		Name:      "Islamabad",
	},
	"ISTANBUL": {
		Longitude: 28.998413,
		Latitude:  41.055908,
		Name:      "Istanbul",
	},
	"JERUSALEM": {
		Longitude: 35.211182,
		Latitude:  31.761475,
		Name:      "Jerusalem",
	},
	"JOHANNESBURG": {
		Longitude: 28.048096,
		Latitude:  -26.141968,
		Name:      "Johannesburg",
	},
	"KUWAIT CITY": {
		Longitude: 47.977295,
		Latitude:  29.366455,
		Name:      "Kuwait City",
	},
	"LAS VEGAS": {
		Longitude: -115.131226,
		Latitude:  36.172485,
		Name:      "Las Vegas",
	},
	"LONDON": {
		Longitude: -0.115356,
		Latitude:  51.503906,
		Name:      "London",
	},
	"LOS ANGELES": {
		Longitude: -118.240356,
		Latitude:  34.052124,
		Name:      "Los Angeles",
	},
	"LUXEMBOURG": {
		Longitude: 6.124878,
		Latitude:  49.608765,
		Name:      "Luxembourg",
	},
	"MADRID": {
		Longitude: -3.702393,
		Latitude:  40.413208,
		Name:      "Madrid",
	},
	"MEXICO CITY": {
		Longitude: -99.135132,
		Latitude:  19.429321,
		Name:      "Mexico City",
	},
	"MIAMI": {
		Longitude: -80.189209,
		Latitude:  25.768433,
		Name:      "Miami",
	},
	"MILAN": {
		Longitude: 9.184570,
		Latitude:  45.466919,
		Name:      "Milan",
	},
	"MILWAUKEE": {
		Longitude: -87.901611,
		Latitude:  43.033447,
		Name:      "Milwaukee",
	},
	"MINNEAPOLIS": {
		Longitude: -93.262939,
		Latitude:  44.978027,
		Name:      "Minneapolis",
	},
	"MONACO": {
		Longitude: 7.432251,
		Latitude:  43.714600,
		Name:      "Monaco",
	},
	"MOSCOW": {
		Longitude: 37.611694,
		Latitude:  55.766602,
		Name:      "Moscow",
	},
	"MUNICH": {
		Longitude: 11.552124,
		Latitude:  48.131104,
		Name:      "Munich",
	},
	"NEW DELHI": {
		Longitude: 77.195435,
		Latitude:  28.597412,
		Name:      "New Delhi",
	},
	"NEW ORLEANS": {
		Longitude: -90.071411,
		Latitude:  29.954224,
		Name:      "New Orleans",
	},
	"NEW YORK": {
		Longitude: -74.003906,
		Latitude:  40.709839,
		Name:      "New York",
	},
	"OKLAHOMA CITY": {
		Longitude: -97.514648,
		Latitude:  35.463867,
		Name:      "Oklahoma City",
	},
	"PANAMA CITY": {
		Longitude: -79.530029,
		Latitude:  8.964844,
		Name:      "Panama City",
	},
	"PARIS": {
		Longitude: 2.345581,
		Latitude:  48.850708,
		Name:      "Paris",
	},
	"PHILADELPHIA": {
		Longitude: -75.162964,
		Latitude:  39.951782,
		Name:      "Philadelphia",
	},
	"PHOENIX": {
		Longitude: -112.071533,
		Latitude:  33.447876,
		Name:      "Phoenix",
	},
	"PITTSBURGH": {
		Longitude: -79.991455,
		Latitude:  40.435181,
		Name:      "Pittsburgh",
	},
	"PRAGUE": {
		Longitude: 14.430542,
		Latitude:  50.070190,
		Name:      "Prague",
	},
	"RIO DE JANEIRO": {
		Longitude: -43.231201,
		Latitude:  -22.895508,
		Name:      "Rio de Janeiro",
	},
	"ROME": {
		Longitude: 12.485962,
		Latitude:  41.890869,
		Name:      "Rome",
	},
	"SALT LAKE CITY": {
		Longitude: -111.890259,
		Latitude:  40.759277,
		Name:      "Salt Lake City",
	},
	"SAN ANTONIO": {
		Longitude: -98.492432,
		Latitude:  29.421387,
		Name:      "San Antonio",
	},
	"SAN DIEGO": {
		Longitude: -67.862549,
		Latitude:  10.250244,
		Name:      "San Diego",
	},
	"SAN FRANCISCO": {
		Longitude: -122.415161,
		Latitude:  37.770996,
		Name:      "San Francisco",
	},
	"SAN MARINO": {
		Longitude: 12.431030,
		Latitude:  43.928833,
		Name:      "San Marino",
	},
	"SEATTLE": {
		Longitude: -122.327271,
		Latitude:  47.603760,
		Name:      "Seattle",
	},
	"SHANGHAI": {
		Longitude: 121.470337,
		Latitude:  31.245117,
		Name:      "Shanghai",
	},
	"SINGAPORE": {
		Longitude: 103.853760,
		Latitude:  1.290894,
		Name:      "Singapore",
	},
	"ST. LOUIS": {
		Longitude: -90.197754,
		Latitude:  38.622437,
		Name:      "St. Louis",
	},
	"STOCKHOLM": {
		Longitude: 18.072510,
		Latitude:  59.282227,
		Name:      "Stockholm",
	},
	"SYDNEY": {
		Longitude: 151.237793,
		Latitude:  -33.887329,
		Name:      "Sydney",
	},
	"TOKYO": {
		Longitude: 139.762573,
		Latitude:  35.683594,
		Name:      "Tokyo",
	},
	"TORONTO": {
		Longitude: -79.414673,
		Latitude:  43.698120,
		Name:      "Toronto",
	},
	"VATICAN CITY": {
		Longitude: 12.453483,
		Latitude:  41.903512,
		Name:      "Vatican City",
	},
	"VIENNA": {
		Longitude: 16.369629,
		Latitude:  48.202515,
		Name:      "Vienna",
	},
	"WASHINGTON": {
		Longitude: -77.036133,
		Latitude:  38.891602,
		Name:      "Washington D.C.",
	},
	"MACAU": {
		Longitude: 113.5986,
		Latitude:  22.21435,
		Name:      "Macao",
	},
	"MONTREAL": {
		Longitude: -73.646850,
		Latitude:  45.516357,
		Name:      "Montreal",
	},
	"QUEBEC CITY": {
		Longitude: -71.20788,
		Latitude:  46.8017,
		Name:      "Quebec City",
	},
	"SAO PAULO": {
		Longitude: -46.614990,
		Latitude:  -23.53271,
		Name:      "Sao Paulo",
	},
	"ZURICH": {
		Longitude: 8.5363769,
		Latitude:  47.3895263,
		Name:      "Zurich",
	},
}
