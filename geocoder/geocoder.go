package geocoder

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/joho/godotenv"
	"io/ioutil"
	"net/http"
	"os"
)

type Coordinates struct {
	Latitude  float64
	Longitude float64
}

func GetCityCoordinates(city string) (*Coordinates, error) {
	geoObj, err := fetchCity(city)
	if err != nil {
		return &Coordinates{}, err
	}

	coords, err := geoObj.Point.toCoords()

	return &coords, nil
}

func fetchCity(city string) (*GeoObject, error) {
	err := godotenv.Load()
	if err != nil {
		return &GeoObject{}, errors.New("Failed to load environment variables. Make sure .env file exists")
	}

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
		return &GeoObject{}, errors.New("Failed to parse response as JSON data")
	}

	var featureMember []FeatureMember
	featureMember = response.Response.GeoObjectCollection.FeatureMember

	if len(featureMember) == 0 {
		return &GeoObject{}, errors.New("City not found")
	}

	return &featureMember[0].GeoObject, nil
}
