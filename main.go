package main

import (
	"./weather"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
)

type CurrentWeatherRequestBody struct {
	City string `json:"city"`
}

func main() {
	fmt.Println()

	err := weather.Init()
	if err != nil {
		log.Fatalln(err)
	}

	r := mux.NewRouter()
	r.HandleFunc("/api/weather/current/{city}", getCurrentWeatherHandler)
	r.HandleFunc("/api/weather/current", getCurrentWeatherHandler)
	r.NotFoundHandler = http.HandlerFunc(notFound)

	fs := http.FileServer(http.Dir("./public"))
	r.PathPrefix("/").Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			http.ServeFile(w, r, "public/index.html")
			return
		}

		fs.ServeHTTP(w, r)
	})).Methods("GET")

	log.Println(http.ListenAndServe(":8080", r))
}

func showIndex(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "public/index.html")
}

func getCurrentWeatherHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var city string

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
		http.ServeFile(w, r, "public/400.txt")
		return
	}

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
		Temperature string `json:"temperature"`
		FeelsLike   string `json:"feels_like"`
	}{
		Temperature: prettifyTemperature(weatherInfo.Current.TempC),
		FeelsLike:   prettifyTemperature(weatherInfo.Current.FeelsLikeC),
	}

	err = json.NewEncoder(w).Encode(&filteredWeather)
	if err != nil {
		log.Fatalln(err)
	}
}

func prettifyTemperature(t float64) string {
	tStr := fmt.Sprintf("%.1f", t)

	if tStr[len(tStr)-2:] == ".0" {
		tStr = tStr[:2]
	}

	if t == 0.0 {
		tStr = "0"
	}

	if t > 0 {
		tStr = "+" + tStr
	}

	return tStr
}

func notFound(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "public/404.txt")
}
