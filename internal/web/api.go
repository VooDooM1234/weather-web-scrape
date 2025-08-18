package web

import (
	"encoding/json"
	"net/http"

	"weather-scraper.com/internal/fetch"
)

func ServeWeatherJSON(w http.ResponseWriter, r *http.Request, data fetch.WeatherResponse) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}
