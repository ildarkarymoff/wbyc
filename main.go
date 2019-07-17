package main

import (
	"./weather"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

func main() {
	fmt.Println()

	err := weather.Init() // Initializing weather module to work with Redis
	if err != nil {
		log.Fatalln(err)
	}

	r := mux.NewRouter()
	r.HandleFunc("/api/weather/current/{city}", handleGetCurrentWeather)
	r.HandleFunc("/api/weather/current", handleGetCurrentWeather)
	r.NotFoundHandler = http.HandlerFunc(handleNotFound)

	err = http.ListenAndServe(":8080", r)
	if err != nil {
		log.Fatalln(err)
	}
}

func handleGetCurrentWeather(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var city string

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
		http.ServeFile(w, r, "public/400.txt")
		return
	}

	// Parsing request parameters from query, URI or request body:
	params, err := resolveParameters(r, &body, "city", "icon")
	if err != nil {
		log.Println(err)
		http.ServeFile(w, r, "public/400.txt")
		return
	}

	city = params["city"]

	log.Printf("%s (%s), %s", r.Method, r.Header.Get("Content-Type"), city)
	getCurrentWeather(w, r, city)
}

func resolveParameters(r *http.Request, body *[]byte, keys ...string) (map[string]string, error) {
	uriParams := mux.Vars(r)
	queryParams := r.URL.Query()
	resolved := make(map[string]string)
	var requestBody map[string]interface{}

	if len(*body) > 0 {
		err := json.Unmarshal(*body, &requestBody)
		if err != nil {
			return map[string]string{}, err
		}
	}

	for _, key := range keys {
		if val, exists := uriParams[key]; exists {
			resolved[key] = val
		} else if len(queryParams[key]) > 0 {
			resolved[key] = queryParams[key][0]
		} else if val, exists := requestBody[key]; exists {
			resolved[key] = fmt.Sprintf("%v", val)
		}
	}

	return resolved, nil
}

func getCurrentWeather(w http.ResponseWriter, r *http.Request, city string) {
	weatherInfo, err := weather.GetCurrent(city)
	if err != nil {
		log.Printf("Failed to get weather: %v", err)
		w.WriteHeader(http.StatusNotFound)
		http.ServeFile(w, r, "public/404.txt")
		return
	}

	filteredWeather := struct {
		City        string `json:"city"`
		Temperature string `json:"temperature"`
		FeelsLike   string `json:"feels_like"`
	}{
		City:        weatherInfo.City,
		Temperature: prettifyTemperature(weatherInfo.Weather.Current.TempC),
		FeelsLike:   prettifyTemperature(weatherInfo.Weather.Current.FeelsLikeC),
	}

	err = json.NewEncoder(w).Encode(&filteredWeather)
	if err != nil {
		log.Fatalln(err)
	}
}

func prettifyTemperature(t float64) string {
	tStr := strconv.FormatFloat(t, 'f', -1, 64)

	if t > 0 {
		tStr = "+" + tStr
	}

	return tStr
}

func handleNotFound(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "public/404.txt")
}
