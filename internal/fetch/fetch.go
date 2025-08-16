package fetch

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

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

	TempUnit string
}

type WeatherCurrentView struct {
	LastUpdated string
	TempC       float64
	TempF       float64
	IsDay       int
	WindMph     float64
	WindKph     float64
	WindDegree  int
	WindDir     string
	PressureMb  float64
	PressureIn  float64
	PrecipMm    float64
	PrecipIn    float64
	Humidity    int
	Cloud       int
	FeelsLikeC  float64
	FeelsLikeF  float64
	WindChillC  float64
	WindChillF  float64
	HeatIndexC  float64
	HeatIndexF  float64
	DewPointC   float64
	DewPointF   float64
	VisKm       float64
	VisMiles    float64
	UV          float64
	GustMph     float64
	GustKph     float64
}

type FetchData struct {
	Scheme   string
	Host     string
	Port     string
	EndPoint string

	Data WeatherResponse
}

func NewWeatherCurrentView(data WeatherResponse) WeatherCurrentView {
	return WeatherCurrentView{
		LastUpdated: data.Current.LastUpdated,
		TempC:       data.Current.TempC,
		TempF:       data.Current.TempF,
		IsDay:       data.Current.IsDay,
		WindMph:     data.Current.WindMph,
		WindKph:     data.Current.WindKph,
		WindDegree:  data.Current.WindDegree,
		WindDir:     data.Current.WindDir,
		PressureMb:  data.Current.PressureMb,
		PressureIn:  data.Current.PressureIn,
		PrecipMm:    data.Current.PrecipMm,
		PrecipIn:    data.Current.PrecipIn,
		Humidity:    data.Current.Humidity,
		Cloud:       data.Current.Cloud,
		FeelsLikeC:  data.Current.FeelsLikeC,
		FeelsLikeF:  data.Current.FeelsLikeF,
		WindChillC:  data.Current.WindChillC,
		WindChillF:  data.Current.WindChillF,
		HeatIndexC:  data.Current.HeatIndexC,
		HeatIndexF:  data.Current.HeatIndexF,
		DewPointC:   data.Current.DewPointC,
		DewPointF:   data.Current.DewPointF,
		VisKm:       data.Current.VisKm,
		VisMiles:    data.Current.VisMiles,
		UV:          data.Current.UV,
		GustMph:     data.Current.GustMph,
		GustKph:     data.Current.GustKph,
	}
}

func (f FetchData) CreateWeatherLabels() map[string]string {
	return map[string]string{

		"Name":           "Name",
		"Region":         "Region",
		"Country":        "Country",
		"Lat":            "Latitude",
		"Lon":            "Longitude",
		"TzID":           "Timezone ID",
		"LocaltimeEpoch": "Local Epoch Time",
		"Localtime":      "Local Time",

		"LastUpdated": "Last Updated",
		"TempC":       "Temperature (°C)",
		"TempF":       "Temperature (°F)",
		"IsDay":       "Daytime (1=Yes, 0=No)",
		"WindMph":     "Wind (mph)",
		"WindKph":     "Wind (kph)",
		"WindDegree":  "Wind Direction (°)",
		"WindDir":     "Wind Direction",
		"PressureMb":  "Pressure (mb)",
		"PressureIn":  "Pressure (in)",
		"PrecipMm":    "Precipitation (mm)",
		"PrecipIn":    "Precipitation (in)",
		"Humidity":    "Humidity (%)",
		"Cloud":       "Cloud Cover (%)",
		"FeelsLikeC":  "Feels Like (°C)",
		"FeelsLikeF":  "Feels Like (°F)",
		"WindChillC":  "Wind Chill (°C)",
		"WindChillF":  "Wind Chill (°F)",
		"HeatIndexC":  "Heat Index (°C)",
		"HeatIndexF":  "Heat Index (°F)",
		"DewPointC":   "Dew Point (°C)",
		"DewPointF":   "Dew Point (°F)",
		"VisKm":       "Visibility (km)",
		"VisMiles":    "Visibility (miles)",
		"UV":          "UV Index",
		"GustMph":     "Wind Gust (mph)",
		"GustKph":     "Wind Gust (kph)",
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
