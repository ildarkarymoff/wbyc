package weather

import (
	"./apixu"
	"./geocoder"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis"
	"log"
	"os"
	"time"
)

var redisClient *redis.Client

// Init creates new Redis client (required for storing weather cache)
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

// GetCurrent fetches city coordinates by it's name,
// then fetches current weather data from Redis (if
// corresponding coordinates-key exists) or directly
// from Apixu Weather.
func GetCurrent(city string) (*apixu.Weather, error) {
	var weather *apixu.Weather

	coordinates, err := geocoder.GetCoordinates(city)
	if err != nil {
		return &apixu.Weather{}, err
	}

	if existsInRedis(coordinates) {
		weather, err = getWeatherRow(coordinates)
	} else {
		weather, err = apixu.GetCurrentWeather(coordinates)
		if err != nil {
			return &apixu.Weather{}, err
		}

		err = setWeatherRow(coordinates, weather)
	}

	if err != nil {
		return &apixu.Weather{}, err
	}

	return weather, err
}

func existsInRedis(coordinates *geocoder.Coordinates) bool {
	exists, err := redisClient.Exists(makeKeyFromCoords(coordinates)).Result()
	if err != nil {
		log.Printf("[existsInRedis] %v", err)
		return false
	}

	return !(exists == 0)
}

func makeKeyFromCoords(coordinates *geocoder.Coordinates) string {
	return fmt.Sprintf("%f,%f",
		coordinates.Latitude,
		coordinates.Longitude)
}

func setWeatherRow(coordinates *geocoder.Coordinates, weather *apixu.Weather) error {
	serialized, _ := json.Marshal(weather)
	err := redisClient.Set(
		makeKeyFromCoords(coordinates),
		string(serialized),
		60*60*time.Second, // Weather "cache" lives for 1 hour
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
