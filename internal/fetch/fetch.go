package fetch

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

// API weather response structure
type WeatherResponse struct {
	Location struct {
		Name           string  `json:"name"`
		Region         string  `json:"region"`
		Country        string  `json:"country"`
		Lat            float64 `json:"lat"`
		Lon            float64 `json:"lon"`
		TzID           string  `json:"tz_id"`
		LocaltimeEpoch int64   `json:"localtime_epoch"`
		Localtime      string  `json:"localtime"`
	} `json:"location"`

	Current struct {
		LastUpdatedEpoch int64   `json:"last_updated_epoch"`
		LastUpdated      string  `json:"last_updated"`
		TempC            float64 `json:"temp_c"`
		TempF            float64 `json:"temp_f"`
		IsDay            int     `json:"is_day"`
		Condition        struct {
			Text string `json:"text"`
			Icon string `json:"icon"`
			Code int    `json:"code"`
		} `json:"condition"`
		WindMph    float64 `json:"wind_mph"`
		WindKph    float64 `json:"wind_kph"`
		WindDegree int     `json:"wind_degree"`
		WindDir    string  `json:"wind_dir"`
		PressureMb float64 `json:"pressure_mb"`
		PressureIn float64 `json:"pressure_in"`
		PrecipMm   float64 `json:"precip_mm"`
		PrecipIn   float64 `json:"precip_in"`
		Humidity   int     `json:"humidity"`
		Cloud      int     `json:"cloud"`
		FeelsLikeC float64 `json:"feelslike_c"`
		FeelsLikeF float64 `json:"feelslike_f"`
		WindChillC float64 `json:"windchill_c"`
		WindChillF float64 `json:"windchill_f"`
		HeatIndexC float64 `json:"heatindex_c"`
		HeatIndexF float64 `json:"heatindex_f"`
		DewPointC  float64 `json:"dewpoint_c"`
		DewPointF  float64 `json:"dewpoint_f"`
		VisKm      float64 `json:"vis_km"`
		VisMiles   float64 `json:"vis_miles"`
		UV         float64 `json:"uv"`
		GustMph    float64 `json:"gust_mph"`
		GustKph    float64 `json:"gust_kph"`
	} `json:"current"`
}

// Struct for rendering weather data view in html template
type WeatherView struct {
	Data     WeatherResponse
	TempUnit string
}

type FetchData struct {
	Scheme   string
	Host     string
	Port     string
	EndPoint string

	Data WeatherResponse
}

type WeatherLocation struct {
	Name           string  `json:"name"`
	Region         string  `json:"region"`
	Country        string  `json:"country"`
	Lat            float64 `json:"lat"`
	Lon            float64 `json:"lon"`
	TzID           string  `json:"tz_id"`
	LocaltimeEpoch int64   `json:"localtime_epoch"`
	Localtime      string  `json:"localtime"`
}

type WeatherCurrent struct {
	LastUpdatedEpoch int64   `json:"last_updated_epoch"`
	LastUpdated      string  `json:"last_updated"`
	TempC            float64 `json:"temp_c"`
	TempF            float64 `json:"temp_f"`
	IsDay            int     `json:"is_day"`

	ConditionText string `json:"condition_text"`
	ConditionIcon string `json:"condition_icon"`
	ConditionCode int    `json:"condition_code"`

	WindMph    float64 `json:"wind_mph"`
	WindKph    float64 `json:"wind_kph"`
	WindDegree int     `json:"wind_degree"`
	WindDir    string  `json:"wind_dir"`

	PressureMb float64 `json:"pressure_mb"`
	PressureIn float64 `json:"pressure_in"`
	PrecipMm   float64 `json:"precip_mm"`
	PrecipIn   float64 `json:"precip_in"`

	Humidity   int     `json:"humidity"`
	Cloud      int     `json:"cloud"`
	FeelsLikeC float64 `json:"feelslike_c"`
	FeelsLikeF float64 `json:"feelslike_f"`
	WindChillC float64 `json:"windchill_c"`
	WindChillF float64 `json:"windchill_f"`
	HeatIndexC float64 `json:"heatindex_c"`
	HeatIndexF float64 `json:"heatindex_f"`
	DewPointC  float64 `json:"dewpoint_c"`
	DewPointF  float64 `json:"dewpoint_f"`

	VisKm    float64 `json:"vis_km"`
	VisMiles float64 `json:"vis_miles"`
	UV       float64 `json:"uv"`
	GustMph  float64 `json:"gust_mph"`
	GustKph  float64 `json:"gust_kph"`
}

type WeatherResponseFlat struct {
	Location WeatherLocation `json:"location"`
	Current  WeatherCurrent  `json:"current"`
}

func FlattenWeather(orig *WeatherResponse) WeatherResponseFlat {
	return WeatherResponseFlat{
		Location: WeatherLocation{
			Name:           orig.Location.Name,
			Region:         orig.Location.Region,
			Country:        orig.Location.Country,
			Lat:            orig.Location.Lat,
			Lon:            orig.Location.Lon,
			TzID:           orig.Location.TzID,
			LocaltimeEpoch: orig.Location.LocaltimeEpoch,
			Localtime:      orig.Location.Localtime,
		},
		Current: WeatherCurrent{
			LastUpdatedEpoch: orig.Current.LastUpdatedEpoch,
			LastUpdated:      orig.Current.LastUpdated,
			TempC:            orig.Current.TempC,
			TempF:            orig.Current.TempF,
			IsDay:            orig.Current.IsDay,

			ConditionText: orig.Current.Condition.Text,
			// ConditionIcon: orig.Current.Condition.Icon,
			ConditionCode: orig.Current.Condition.Code,

			WindMph:    orig.Current.WindMph,
			WindKph:    orig.Current.WindKph,
			WindDegree: orig.Current.WindDegree,
			WindDir:    orig.Current.WindDir,

			PressureMb: orig.Current.PressureMb,
			PressureIn: orig.Current.PressureIn,
			PrecipMm:   orig.Current.PrecipMm,
			PrecipIn:   orig.Current.PrecipIn,

			Humidity:   orig.Current.Humidity,
			Cloud:      orig.Current.Cloud,
			FeelsLikeC: orig.Current.FeelsLikeC,
			FeelsLikeF: orig.Current.FeelsLikeF,
			WindChillC: orig.Current.WindChillC,
			WindChillF: orig.Current.WindChillF,
			HeatIndexC: orig.Current.HeatIndexC,
			HeatIndexF: orig.Current.HeatIndexF,
			DewPointC:  orig.Current.DewPointC,
			DewPointF:  orig.Current.DewPointF,

			VisKm:    orig.Current.VisKm,
			VisMiles: orig.Current.VisMiles,
			UV:       orig.Current.UV,
			GustMph:  orig.Current.GustMph,
			GustKph:  orig.Current.GustKph,
		},
	}
}

func WeatherLabels() map[string]string {
	return map[string]string{
		"last_updated_epoch": "Last Updated (Epoch)",
		"last_updated":       "Last Updated",
		"temp_c":             "Temperature (°C)",
		"temp_f":             "Temperature (°F)",
		"feelslike_c":        "Feels Like (°C)",
		"feelslike_f":        "Feels Like (°F)",
		"humidity":           "Humidity (%)",
		"wind_mph":           "Wind Speed (mph)",
		"wind_kph":           "Wind Speed (kph)",
		"wind_degree":        "Wind Direction (°)",
		"wind_dir":           "Wind Direction",
		"pressure_mb":        "Pressure (mb)",
		"pressure_in":        "Pressure (inHg)",
		"precip_mm":          "Precipitation (mm)",
		"precip_in":          "Precipitation (in)",
		"uv":                 "UV Index",
		"gust_mph":           "Gust Speed (mph)",
		"gust_kph":           "Gust Speed (kph)",
		"condition_text":     "Condition Text",
		// "condition_icon":     "Condition Icon",
		"condition_code": "Condition Code",
		"windchill_c":    "Wind Chill (°C)",
		"windchill_f":    "Wind Chill (°F)",
		"dewpoint_c":     "Dew Point (°C)",
		"dewpoint_f":     "Dew Point (°F)",
		"cloud":          "Cloud Cover (%)",
		"vis_km":         "Visibility (km)",
		"vis_miles":      "Visibility (miles)",
		"is_day":         "Daytime Flag",
	}
}

func NewFetch(host, endPoint string) *FetchData {
	return &FetchData{
		Scheme:   "https",
		Host:     host,
		EndPoint: endPoint,
	}
}

func (f *FetchData) FetchWeather(apiKey, qParam string) WeatherResponse {
	var data WeatherResponse

	url := fmt.Sprintf("%s://%s%s?key=%s&q=%s", f.Scheme, f.Host, f.EndPoint, apiKey, qParam)
	fmt.Println("Fetching:", url)

	resp, err := http.Get(url)
	if err != nil {
		log.Fatalf("HTTP request error: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error reading response body: %v", err)
	}

	if resp.StatusCode > 299 {
		log.Fatalf("Response failed with status %d\nBody: %s\n", resp.StatusCode, string(body))
	}

	if err := json.Unmarshal(body, &data); err != nil {
		log.Fatalf("JSON decode error: %v\nBody: %s", err, string(body))
	}

	var prettyJSON bytes.Buffer
	error := json.Indent(&prettyJSON, body, "", "\t")
	if error != nil {
		log.Println("JSON parse error: ", error)
	}

	log.Println("Weather Data:", prettyJSON.String())

	return data
}
