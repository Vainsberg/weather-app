package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
	type requestweather struct {
		Wind        float64 `json:"wind_speed_10m"`
		Temperature float64 `json:"temperature_2m"`
	}

	query := r.URL.Query()
	latitudeText := query.Get("latitude")
	longitudeText := query.Get("longitude")

	if latitudeText == "" || longitudeText == "" {
		http.Error(w, "Latitude and longitude are required parameters", http.StatusBadRequest)
		return
	}
	weatherAPIURL := "https://api.open-meteo.com/v1/forecast?latitude=" + latitudeText + "&longitude=" + longitudeText + "&current=temperature_2m,wind_speed_10m"

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
	var requestWeather requestweather
	err = json.Unmarshal(data, &requestWeather)
	if err != nil {
		log.Println("Error parsing weather data:", err)
		http.Error(w, "Error parsing weather data", http.StatusInternalServerError)
		return
	}
	responseText := fmt.Sprintf("Добрый день! Сегодня температура %0.1f градусов, скорость ветра %0.1f м/с.", requestWeather.Temperature, requestWeather.Wind)

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(responseText))

}

func main() {
	http.HandleFunc("/get_weather", handler)

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Ошибка запуска сервера:", err)
	}
}
