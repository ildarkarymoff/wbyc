package apixu

import (
	"../geocoder"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

type TempAndFeelsLikeWeather struct {
	Temperature string `json:"temperature"`
	FeelsLike   string `json:"feels_like"`
}

//func GetCurrentWeather(coordinates *geocoder.Coordinates) (*TempAndFeelsLikeWeather, error) {
//	response, err := GetCurrentWeather(coordinates)
//	if err != nil {
//		return &TempAndFeelsLikeWeather{}, err
//	}
//
//	return &TempAndFeelsLikeWeather{
//		Temperature: prettifyTemperature(response.Current.TempC),
//		FeelsLike:   prettifyTemperature(response.Current.FeelsLikeC),
//	}, nil
//}

func GetCurrentWeather(coordinates *geocoder.Coordinates) (*Weather, error) {
	apiKey := os.Getenv("APIXU_API_KEY")
	urlFormat := "http://api.apixu.com/v1/current.json?key=%s&q=%f,%f"

	// Idk why Apixu waits for {lon, lat} pair instead of {lat, lon} ¯\_(ツ)_/¯
	url := fmt.Sprintf(urlFormat, apiKey, coordinates.Longitude, coordinates.Latitude)

	resp, err := http.Get(url)
	if err != nil {
		return &Weather{}, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return &Weather{}, err
	}

	var response Weather
	err = json.Unmarshal(body, &response)
	if err != nil {
		log.Println(err)
		return &Weather{}, err
	}

	return &response, nil
}
