package geocoder

import (
	"errors"
	"strconv"
	"strings"
)

// Automatically generated from typical JSON response from
// Yandex Geocoder Maps API
type HttpApiResponse struct {
	Response Response `json:"response"`
}
type GeocoderResponseMetaData struct {
	Request string `json:"request"`
	Found   string `json:"found"`
	Results string `json:"results"`
}
type ResponseMetaDataProperty struct {
	GeocoderResponseMetaData GeocoderResponseMetaData `json:"GeocoderResponseMetaData"`
}
type Components struct {
	Kind string `json:"kind"`
	Name string `json:"name"`
}
type Address struct {
	CountryCode string       `json:"country_code"`
	Formatted   string       `json:"formatted"`
	Components  []Components `json:"Components"`
}
type Locality struct {
	LocalityName string `json:"LocalityName"`
}
type SubAdministrativeArea struct {
	SubAdministrativeAreaName string   `json:"SubAdministrativeAreaName"`
	Locality                  Locality `json:"Locality"`
}
type AdministrativeArea struct {
	AdministrativeAreaName string                `json:"AdministrativeAreaName"`
	SubAdministrativeArea  SubAdministrativeArea `json:"SubAdministrativeArea"`
}
type Country struct {
	AddressLine        string             `json:"AddressLine"`
	CountryNameCode    string             `json:"CountryNameCode"`
	CountryName        string             `json:"CountryName"`
	AdministrativeArea AdministrativeArea `json:"AdministrativeArea"`
}
type AddressDetails struct {
	Country Country `json:"Country"`
}
type GeocoderMetaData struct {
	Kind           string         `json:"kind"`
	Text           string         `json:"text"`
	Precision      string         `json:"precision"`
	Address        Address        `json:"Address"`
	AddressDetails AddressDetails `json:"AddressDetails"`
}
type MetaDataProperty struct {
	GeocoderMetaData GeocoderMetaData `json:"GeocoderMetaData"`
}
type Envelope struct {
	LowerCorner string `json:"lowerCorner"`
	UpperCorner string `json:"upperCorner"`
}
type BoundedBy struct {
	Envelope Envelope `json:"Envelope"`
}
type Point struct {
	Pos string `json:"pos"`
}
type GeoObject struct {
	MetaDataProperty MetaDataProperty `json:"metaDataProperty"`
	Description      string           `json:"description"`
	Name             string           `json:"name"`
	BoundedBy        BoundedBy        `json:"boundedBy"`
	Point            Point            `json:"Point"`
}
type FeatureMember struct {
	GeoObject GeoObject `json:"GeoObject"`
}
type GeoObjectCollection struct {
	MetaDataProperty MetaDataProperty `json:"metaDataProperty"`
	FeatureMember    []FeatureMember  `json:"featureMember"`
}
type Response struct {
	GeoObjectCollection GeoObjectCollection `json:"GeoObjectCollection"`
}

func (p *Point) toCoords() (Coordinates, error) {
	splitted := strings.Split(p.Pos, " ")

	if len(splitted) != 2 {
		return Coordinates{}, errors.New("Wrong format of position string")
	}

	latitude, err := strconv.ParseFloat(splitted[0], 64)
	longitude, err := strconv.ParseFloat(splitted[1], 64)

	if err != nil {
		return Coordinates{}, errors.New("Failed to parse coordinates from position string")
	}

	return Coordinates{
		Latitude:  latitude,
		Longitude: longitude,
	}, nil
}
