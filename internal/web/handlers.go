package web

import (
	"fmt"
	"log"
	"net/http"
	"text/template"

	"weather-scraper.com/internal/fetch"
	"weather-scraper.com/utils"
)

// HandlerData holds dependencies for the handlers
type HandlerData struct {
	Client      *fetch.FetchData
	APIKey      string
	DefaultUnit string
}

func NewHandlerData(client *fetch.FetchData, apiKey, defaultUnit string) *HandlerData {
	return &HandlerData{
		Client:      client,
		APIKey:      apiKey,
		DefaultUnit: defaultUnit,
	}
}

// HomePageHandler serves the main dashboard page
func (h *HandlerData) HomePageHandler(currentWeatherData *fetch.WeatherResponse) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		tmpl := template.Must(template.ParseFiles("templates/index.html"))

		view := fetch.WeatherView{
			Data:     currentWeatherData,
			TempUnit: h.DefaultUnit,
		}

		if err := tmpl.Execute(w, view); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

// ExtendedDataTableHandler serves the extended stats table
func (h *HandlerData) ExtendedDataTableHandler(currentWeatherData *fetch.WeatherResponse) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")

		flatWeather := fetch.FlattenWeather(currentWeatherData)
		currentDataMap, _ := utils.StructToMap(flatWeather.Current)

		labels := fetch.WeatherLabels()

		view := map[string]interface{}{
			"Data":   currentDataMap,
			"Labels": labels,
		}

		tmpl := template.Must(template.ParseFiles("templates/sub_table_weather_extended.html"))
		if err := tmpl.Execute(w, view); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

// SetDataUnitsHandler handles unit change (metric/imperial)
func (h *HandlerData) SetDataUnitsHandler(currentWeatherData *fetch.WeatherResponse) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			http.Error(w, "Failed to parse form", http.StatusBadRequest)
			return
		}
		unit := r.PostFormValue("unit")
		log.Printf("HTMX POST received: %s", unit)

		view := fetch.WeatherView{
			Data:     currentWeatherData,
			TempUnit: unit,
		}

		tmpl := template.Must(template.ParseFiles("templates/index.html"))
		if err := tmpl.ExecuteTemplate(w, "dashboard-quick-grid", view); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

// WeatherLocationSearchHandler handles search queries
func (h *HandlerData) WeatherLocationSearchHandler(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("search")
	results, err := h.Client.FetchWeatherSearch(h.APIKey, q)
	if err != nil {
		http.Error(w, "Weather search failed", http.StatusInternalServerError)
		return
	}
	fmt.Printf("search results: %+v\n", results)
	tmpl := template.Must(template.ParseFiles("templates/index.html"))
	if err := tmpl.ExecuteTemplate(w, "search-result", results); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// WeatherLoadSearchLocationHandler loads the full weather for a selected search result
func (h *HandlerData) WeatherLoadSearchLocationHandler(currentWeatherData *fetch.WeatherResponse) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		name := r.PostFormValue("name")
		region := r.PostFormValue("region")
		country := r.PostFormValue("country")

		log.Printf("Search Results: %s, %s, %s", name, region, country)

		// fetch new data and overwrite the pointer's value
		updatedWeather, err := h.Client.FetchWeatherCurrent(h.APIKey, name)
		if err != nil {
			http.Error(w, "Current Weather Get Failed", http.StatusInternalServerError)
			return
		}
		*currentWeatherData = updatedWeather // update the shared struct

		w.Header().Set("Content-Type", "text/html")
		tmpl := template.Must(template.ParseFiles("templates/index.html"))

		view := fetch.WeatherView{
			Data:     currentWeatherData,
			TempUnit: h.DefaultUnit,
		}

		if err := tmpl.Execute(w, view); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}
