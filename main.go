package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"text/template"

	"github.com/joho/godotenv"
	"weather-scraper.com/internal/fetch"
	"weather-scraper.com/utils"
)

func ServeWeatherJSON(w http.ResponseWriter, r *http.Request, data fetch.WeatherResponse) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func main() {
	_ = godotenv.Load()
	apiKeyWeather := os.Getenv("WATHER_API_KEY")
	if apiKeyWeather == "" {
		log.Fatal("Missing WATHER_API_KEY in environment!")
	}
	const defaultLocation string = "melbourne"
	const defaultUnit = "metric"

	fs := http.FileServer(http.Dir("./static/css"))
	http.Handle("/css/", http.StripPrefix("/css/", fs))

	weatherClient := fetch.NewFetchWeather("https", "api.weatherapi.com", 0)
	currentWeatherData, err := weatherClient.FetchWeatherCurrent(apiKeyWeather, defaultLocation)
	if err != nil {
		fmt.Errorf("API Error: %w", err)
	}

	flatWeather := fetch.FlattenWeather(&currentWeatherData)

	// weatherDataMap, _ := utils.StructToMap(weatherData)
	currentDataMap, _ := utils.StructToMap(flatWeather.Current)
	// locationDataMap, _ := utils.StructToMap(flatWeather.Location)

	h1 := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		tmpl := template.Must(template.ParseFiles("templates/index.html"))

		view := fetch.WeatherView{
			Data:     currentWeatherData,
			TempUnit: defaultUnit,
		}

		tmpl.Execute(w, view)
	}

	extendedDataTable := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")

		labels := fetch.WeatherLabels()

		view := map[string]interface{}{
			"Data":   currentDataMap, // Now all keys are top-level
			"Labels": labels,
		}

		tmpl := template.Must(template.ParseFiles("templates/sub_table_weather_extended.html"))
		tmpl.Execute(w, view)
	}

	setDataUnits := func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
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

	weatherLocationSearch := func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query().Get("search")

		results, err := weatherClient.FetchWeatherSearch(apiKeyWeather, q)
		if err != nil {
			http.Error(w, "Weather search failed", http.StatusInternalServerError)
			return
		}
		fmt.Printf("search results: %+v\n", results)

		// Handle results list
		tmpl := template.Must(template.ParseFiles("templates/index.html"))
		if err := tmpl.ExecuteTemplate(w, "search-result", results); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}

	weatherLoadSearchLocation := func(w http.ResponseWriter, r *http.Request) {
		name := r.PostFormValue("name")
		region := r.PostFormValue("region")
		country := r.PostFormValue("country")

		log.Printf("Search Results: %s, %s, %s", name, region, country)

		currentWeatherData, err := weatherClient.FetchWeatherCurrent(apiKeyWeather, name)
		if err != nil {
			http.Error(w, "Current Weather Get Failed", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "text/html")
		tmpl := template.Must(template.ParseFiles("templates/index.html"))

		view := fetch.WeatherView{
			Data:     currentWeatherData,
			TempUnit: defaultUnit,
		}

		if err := tmpl.Execute(w, view); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

	}

	http.HandleFunc("/", h1)
	http.HandleFunc("/weather/table/", extendedDataTable)
	http.HandleFunc("/set-units/", setDataUnits)
	http.HandleFunc("/weather/search/", weatherLocationSearch)
	http.HandleFunc("/weather/load/", weatherLoadSearchLocation)

	http.HandleFunc("/api/weather", func(w http.ResponseWriter, r *http.Request) {
		ServeWeatherJSON(w, r, currentWeatherData)
	})

	http.HandleFunc("/api/weather-extended-table", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
	})

	log.Println("Starting server on 0.0.0.0:8080")
	log.Fatal(http.ListenAndServe("0.0.0.0:8080", nil))
}
