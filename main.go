package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"weather-scraper.com/internal/fetch"
	"weather-scraper.com/internal/web"
)

func main() {

	const defaultLocation string = "melbourne"
	const defaultUnit = "metric"

	_ = godotenv.Load()
	apiKeyWeather := os.Getenv("WATHER_API_KEY")
	if apiKeyWeather == "" {
		log.Fatal("Missing WATHER_API_KEY in environment!")
	}

	weatherClient := fetch.NewFetchData("https", "api.weatherapi.com", 0)
	currentWeatherData, err := weatherClient.FetchWeatherCurrent(apiKeyWeather, defaultLocation)
	if err != nil {
		fmt.Errorf("API Error: %w", err)
	}

	http.HandleFunc("/api/weather", func(w http.ResponseWriter, r *http.Request) {
		web.ServeWeatherJSON(w, r, currentWeatherData)
	})

	fs := http.FileServer(http.Dir("./static/css"))
	http.Handle("/css/", http.StripPrefix("/css/", fs))

	handler := web.NewHandlerData(weatherClient, apiKeyWeather, defaultUnit)

	http.HandleFunc("/", handler.HomePageHandler(&currentWeatherData))
	http.HandleFunc("/weather/table/", handler.ExtendedDataTableHandler(&currentWeatherData))
	http.HandleFunc("/set-units/", handler.SetDataUnitsHandler(&currentWeatherData))
	http.HandleFunc("/weather/search/", handler.WeatherLocationSearchHandler)
	http.HandleFunc("/weather/load/", handler.WeatherLoadSearchLocationHandler(&currentWeatherData))

	log.Println("Starting server on 0.0.0.0:8080")
	log.Fatal(http.ListenAndServe("0.0.0.0:8080", nil))
}
