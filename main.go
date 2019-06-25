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
	log.Println(http.ListenAndServe(":8080", r))
}

func getCurrentWeather(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	city := params["city"]

	weatherInfo, err := weather.GetCurrentWeather(city)
	if err != nil {
		log.Fatalln(err)
	}

	err = json.NewEncoder(w).Encode(&weatherInfo)
	if err != nil {
		log.Fatalln(err)
	}
}
