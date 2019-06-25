package main

import (
	"./weather"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func main() {
	fmt.Println()

	err := weather.Init()
	if err != nil {
		log.Fatalln(err)
	}

	r := mux.NewRouter()
	r.HandleFunc("/api/weather/current/{city}", getCurrentWeather)
	r.NotFoundHandler = http.HandlerFunc(notFound)
	log.Println(http.ListenAndServe(":8080", r))
}

func getCurrentWeather(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	city := params["city"]

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
