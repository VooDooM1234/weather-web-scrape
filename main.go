package main

import (
	"encoding/json"
	"io"
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

func main() {
	_ = godotenv.Load()
	apiKeyWeather := os.Getenv("WATHER_API_KEY")
	if apiKeyWeather == "" {
		log.Fatal("Missing WATHER_API_KEY in environment!")
	}

	fs := http.FileServer(http.Dir("./static/css"))
	http.Handle("/css/", http.StripPrefix("/css/", fs))

	weatherData := DebugWeatherData()

	http.HandleFunc("/api/weather", func(w http.ResponseWriter, r *http.Request) {
		ServeWeatherJSON(w, r, weatherData)
	})

	http.HandleFunc("/api/weather-extended-table", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
	})

	h1 := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		tmpl := template.Must(template.ParseFiles("templates/index.html"))

		tmpl.Execute(w, weatherData)
	}

	extendedDataTable := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		labels, _ := utils.StructToMap(fetch.FetchData{}.CreateWeatherLabels())
		viewData, _ := utils.StructToMap(fetch.NewWeatherCurrentView(weatherData))

		tmpl := template.Must(template.ParseFiles("templates/sub_table_weather_extended.html"))
		if err := tmpl.Execute(w, map[string]interface{}{
			"Weather": viewData,
			"Labels":  labels,
		}); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}

	setDataUnits := func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Failed to parse form", http.StatusBadRequest)
			return
		}
		unit := r.PostFormValue("unit")
		log.Printf("HTMX POST received: %s", unit)

		weatherData.TempUnit = unit // set the unit

		tmpl := template.Must(template.ParseFiles("templates/index.html"))
		if err := tmpl.ExecuteTemplate(w, "dashboard-quick-grid", weatherData); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}

	http.HandleFunc("/weather/table", extendedDataTable)
	http.HandleFunc("/set-units/", setDataUnits)
	http.HandleFunc("/", h1)

	log.Println("Starting server on 0.0.0.0:8080")
	log.Fatal(http.ListenAndServe("0.0.0.0:8080", nil))
}
