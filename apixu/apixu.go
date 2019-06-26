package apixu

import (
	"../geocoder"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

type Weather struct {
	Temperature string `json:"temperature"`
	FeelsLike   string `json:"feels_like"`
}

func GetCurrentWeather(coordinates *geocoder.Coordinates) (*Weather, error) {
	current, err := fetchCurrentWeather(coordinates)
	if err != nil {
		return &Weather{}, err
	}

	return &Weather{
		Temperature: fmt.Sprintf("%.2f", current.TempC),
		FeelsLike:   fmt.Sprintf("%.2f", current.FeelslikeC),
	}, nil
}

func fetchCurrentWeather(coordinates *geocoder.Coordinates) (*Current, error) {
	apiKey := os.Getenv("APIXU_API_KEY")
	urlFormat := "http://api.apixu.com/v1/current.json?key=%s&q=%f,%f"

	// Idk why Apixu waits {lon, lat} pair instead of {lat, lon} ¯\_(ツ)_/¯
	url := fmt.Sprintf(urlFormat, apiKey, coordinates.Longitude, coordinates.Latitude)

	resp, err := http.Get(url)
	if err != nil {
		return &Current{}, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return &Current{}, err
	}

	var response HttpApiResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		log.Println(err)
		return &Current{}, errors.New("Failed to parse response as JSON data")
	}

	return &response.Current, nil
}
