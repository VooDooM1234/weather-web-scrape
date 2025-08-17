package main

//add weather stuff in here in clean up

// import (
// 	"net/http"
// 	"text/template"
// )

// func weatherLocationSearch(w http.ResponseWriter, r *http.Request) {
// 	query := r.URL.Query().Get("search")
// 	locations, err := weatherClient.FetchWeatherSearch(apiKeyWeather, query)
// 	if err != nil {
// 		http.Error(w, "Failed to fetch locations", http.StatusInternalServerError)
// 		return
// 	}

// 	tmpl := template.Must(template.ParseFiles("templates/index.html"))
// 	tmpl.ExecuteTemplate(w, "search-result", locations)
// }
