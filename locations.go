package main

import (
	"math"
	"unicode/utf16"
)

type Location struct {
	TextOffset   uint32
	Latitude     int16
	Longitude    int16
	CountryCode  uint8
	RegionCode   uint8
	LocationCode uint16
	Zoom         uint8
	_            [3]byte
}

const float64EqualityThreshold = 1e-9

func floatCompare(a, b float64) bool {
	return math.Abs(a-b) <= float64EqualityThreshold
}

func CoordinateEncode(value float64) int16 {
	value /= 0.0054931640625
	return int16(value)
}

func (n *News) MakeLocationTable() {
	n.Header.LocationTableOffset = n.GetCurrentSize()

	for _, location := range n.locations {
		n.Locations = append(n.Locations, Location{
			TextOffset:   0,
			Latitude:     CoordinateEncode(location.Latitude),
			Longitude:    CoordinateEncode(location.Longitude),
			CountryCode:  0,
			RegionCode:   0,
			LocationCode: 0,
			Zoom:         6,
		})
	}

	for i, location := range n.locations {
		n.Locations[i].TextOffset = n.GetCurrentSize()
		encoded := utf16.Encode([]rune(location.Name))
		n.LocationText = append(n.LocationText, encoded...)
		n.LocationText = append(n.LocationText, 0)
		for n.GetCurrentSize()%4 != 0 {
			n.LocationText = append(n.LocationText, 0)
		}
	}

	n.Header.NumberOfLocations = uint32(len(n.Locations))
}
