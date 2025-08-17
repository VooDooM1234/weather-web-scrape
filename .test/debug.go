package main

import (
	"encoding/json"
	"io"
	"log"
	"os"

	"weather-scraper.com/internal/fetch"
)

func DebugWeatherData() fetch.WeatherResponse {
	jsonFile, err := os.Open(".test/test_weather.json")
	if err != nil {
		log.Fatalf("Failed to open JSON file: %v", err)
	}
	defer jsonFile.Close()

	byteValue, err := io.ReadAll(jsonFile)
	if err != nil {
		log.Fatalf("Failed to read JSON file: %v", err)
	}

	var weatherData fetch.WeatherResponse
	err = json.Unmarshal(byteValue, &weatherData)
	if err != nil {
		log.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	return weatherData
}
