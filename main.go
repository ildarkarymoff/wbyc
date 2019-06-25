package main

import (
	"./apixu"
	"./geocoder"
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/api/weather/current/{city}", getCurrentWeather)
	log.Println(http.ListenAndServe(":8080", r))
}

func getCurrentWeather(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	city := params["city"]

	coords, err := geocoder.GetCityCoordinates(city)
	if err != nil {
		log.Fatalln(err)
	}

	weather, err := apixu.GetCurrentWeather(*coords)
	if err != nil {
		log.Fatalln(err)
	}

	err = json.NewEncoder(w).Encode(&weather)
	if err != nil {
		log.Fatalln(err)
	}
}
