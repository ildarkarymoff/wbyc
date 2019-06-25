package weather

import (
	"../apixu"
	"../geocoder"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis"
	"os"
	"time"
)

var redisClient *redis.Client

func Init() error {
	redisClient = redis.NewClient(&redis.Options{
		Addr: os.Getenv("REDIS_HOST") + ":6379",
	})

	_, err := redisClient.Ping().Result()
	if err != nil {
		return err
	}

	return nil
}

func GetCurrentWeather(city string) (*apixu.Weather, error) {
	var weather *apixu.Weather

	coordinates, err := geocoder.GetCityCoordinates(city)
	if err != nil {
		return &apixu.Weather{}, err
	}

	exists, err := redisClient.Exists(makeKeyFromCoords(coordinates)).Result()
	if err != nil {
		return &apixu.Weather{}, err
	}

	if exists == 0 {
		weather, err = apixu.GetCurrentWeather(coordinates)
		_ = setWeatherRow(coordinates, weather)
	} else {
		weather, err = getWeatherRow(coordinates)
	}

	return weather, err
}

func setWeatherRow(coordinates *geocoder.Coordinates, weather *apixu.Weather) error {
	serialized, _ := json.Marshal(weather)
	err := redisClient.Set(
		makeKeyFromCoords(coordinates),
		string(serialized),
		60*60*time.Second,
	)

	return err.Err()
}

func getWeatherRow(coordinates *geocoder.Coordinates) (*apixu.Weather, error) {
	serialized, err := redisClient.Get(makeKeyFromCoords(coordinates)).Result()
	if err != nil {
		return &apixu.Weather{}, nil
	}

	var weather *apixu.Weather
	err = json.Unmarshal([]byte(serialized), &weather)

	return weather, err
}

func makeKeyFromCoords(coordinates *geocoder.Coordinates) string {
	return fmt.Sprintf("%f,%f",
		coordinates.Latitude,
		coordinates.Longitude)
}
