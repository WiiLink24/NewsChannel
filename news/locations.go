package news

type Location struct {
	Longitude float64
	Latitude  float64
	Name      string
}

var CommonLocations = map[string]Location{
	"WASHINGTON": {
		Longitude: -77.0369,
		Latitude:  38.8916015625,
		Name:      "Washington D.C",
	},
	"SEATTLE": {
		Longitude: -122.3272705078125,
		Latitude:  47.603759765625,
		Name:      "Seattle",
	},
}
