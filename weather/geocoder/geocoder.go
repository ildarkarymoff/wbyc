package geocoder

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

type Coordinates struct {
	City      string
	Latitude  float64
	Longitude float64
}

type ErrCityNotFound struct {
	city string
}

func NewErrCityNotFound(city string) *ErrCityNotFound {
	return &ErrCityNotFound{
		city: city,
	}
}

func (e *ErrCityNotFound) Error() string {
	return fmt.Sprintf("City '%s' not found", e.city)
}

func GetCoordinates(city string) (*Coordinates, error) {
	geoObj, err := fetchCity(city)
	if err != nil {
		return &Coordinates{}, err
	}

	coords, err := geoObj.Point.toCoords()
	coords.City = geoObj.Name

	return &coords, nil
}

func fetchCity(city string) (*GeoObject, error) {
	apiKey := os.Getenv("YANDEX_API_KEY")
	urlFormat := "https://geocode-maps.yandex.ru/1.x/?apikey=%s&format=json&geocode=%s"
	url := fmt.Sprintf(urlFormat, apiKey, city)

	resp, err := http.Get(url)
	if err != nil {
		return &GeoObject{}, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return &GeoObject{}, err
	}

	var response HttpApiResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return &GeoObject{}, err
	}

	var featureMember []FeatureMember
	featureMember = response.Response.GeoObjectCollection.FeatureMember

	if len(featureMember) == 0 {
		return &GeoObject{}, NewErrCityNotFound(city)
	}

	return &featureMember[0].GeoObject, nil
}
