package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {

	query := r.URL.Query()
	latitudeText := query.Get("latitude")
	longitudeText := query.Get("longitude")

	if latitudeText == "" || longitudeText == "" {
		http.Error(w, "Latitude and longitude are required parameters", http.StatusBadRequest)
		return
	}
	weatherAPIURL := "https://api.open-meteo.com/v1/forecast?latitude=" + latitudeText + "&longitude=" + longitudeText

	weatherResp, err := http.Get(weatherAPIURL)
	if err != nil {
		log.Println("Error fetching weather data:", err)
		http.Error(w, "Error fetching weather data", http.StatusInternalServerError)
		return
	}
	defer weatherResp.Body.Close()

	data, err := io.ReadAll(weatherResp.Body)
	if err != nil {
		log.Println(err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func main() {
	http.HandleFunc("/get_weather", handler)

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Ошибка запуска сервера:", err)
	}
}
