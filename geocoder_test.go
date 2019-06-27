package main

import (
	. "./weather/geocoder"
	"log"
	"testing"
)

type geoCoderResponseTestPair struct {
	city   string
	coords Coordinates
}

func TestParseCoordinates(t *testing.T) {
	var tests = []geoCoderResponseTestPair{
		{
			city: "Питер",
			coords: Coordinates{
				Latitude:  30.315868,
				Longitude: 59.939095,
			},
		},
		{
			city: "мск",
			coords: Coordinates{
				Latitude:  37.617635,
				Longitude: 55.755814,
			},
		},
		{
			city: "omsk",
			coords: Coordinates{
				Latitude:  73.368212,
				Longitude: 54.989342,
			},
		},
	}

	for _, test := range tests {
		log.Println(test.city)

		coords, err := GetCoordinates(test.city)
		if err != nil {
			t.Errorf("Got error: %v", err)
		}

		if !compareCoords(*coords, test.coords) {
			t.Errorf("Got %v, expected %v", coords, test.coords)
		}
	}
}

func compareCoords(a Coordinates, b Coordinates) bool {
	return a.Latitude == b.Latitude && a.Longitude == b.Longitude
}
