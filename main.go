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
	log.Println(http.ListenAndServe(":8080", r))
}

func getCurrentWeatherHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var city string

	if r.Method == "GET" {
		params := mux.Vars(r)

		if _, exists := params["city"]; !exists {
			cityKeys := r.URL.Query()["city"]
			if len(cityKeys) < 1 {
				http.ServeFile(w, r, "public/404.txt")
				return
			}

			city = cityKeys[0]
		} else {
			city = params["city"]
		}

	} else if r.Method == "POST" {
		var requestBody CurrentWeatherRequestBody

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Println(err)
			http.ServeFile(w, r, "public/400.txt")
			return
		}

		err = json.Unmarshal(body, &requestBody)
		if err != nil {
			log.Println(err)
			http.ServeFile(w, r, "public/400.txt")
			return
		}

		city = requestBody.City
	} else {
		http.ServeFile(w, r, "public/400.txt")
		return
	}

	log.Printf("%s (%s), %s", r.Method, r.Header.Get("Content-Type"), city)

	getCurrentWeather(w, r, city)
}

func getCurrentWeather(w http.ResponseWriter, r *http.Request, city string) {
	weatherInfo, err := weather.GetCurrentWeather(city)
	if err != nil {
		log.Printf("Failed to get weather: %v", err)
		http.ServeFile(w, r, "public/404.txt")
		return
	}

	err = json.NewEncoder(w).Encode(&weatherInfo)
	if err != nil {
		log.Fatalln(err)
	}
}

func notFound(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "public/404.txt")
}
