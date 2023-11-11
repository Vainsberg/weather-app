package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

type responseweather struct {
	Сurrent struct {
		Wind        float64 `json:"wind_speed_10m"`
		Temperature float64 `json:"temperature_2m"`
	} `json:"current"`
}
type Coordinates struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

func handler(w http.ResponseWriter, r *http.Request) {

	var coord Coordinates

	body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Println(http.StatusBadRequest)
		return
	}
	err = json.Unmarshal(body, &coord)
	if err != nil {
		fmt.Println(http.StatusBadRequest)
		return
	}
	row := db.QueryRow("SELECT * FROM text WHERE latitude = $1 AND longitude = $2 AND created_at >= NOW() - INTERVAL 1 HOUR;", coord.Latitude, coord.Longitude)
	responseN := Coordinates{}

	err = row.Scan(&responseN.Latitude, &responseN.Longitude)
	if err != nil {
		log.Fatal(err)
	}

	if responseN.Latitude != 0 || responseN.Longitude != 0 {
		var response responseweather
		responseText := fmt.Sprintf("Добрый день! Сегодня температура %0.1f градусов, скорость ветра %0.1f м/с.", response.Сurrent.Temperature, response.Сurrent.Wind)
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(responseText))
	} else {
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
		var responseWeather responseweather
		err = json.Unmarshal(data, &responseWeather)
		if err != nil {
			log.Println("Error parsing weather data:", err)
			http.Error(w, "Error parsing weather data", http.StatusInternalServerError)
			return
		}
		responseText := fmt.Sprintf("Добрый день! Сегодня температура %0.1f градусов, скорость ветра %0.1f м/с.", responseWeather.Сurrent.Temperature, responseWeather.Сurrent.Wind)

		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(responseText))

	}

}

func main() {
	var err error

	db, err = sql.Open("sqlite3", "store.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS text (
			id INTEGER PRIMARY KEY,
			latitude INT,
			longitude INT
		)
	`)

	if err != nil {
		log.Fatal(err)
	}
	http.HandleFunc("/get_weather", handler)

	errors := http.ListenAndServe(":8080", nil)
	if errors != nil {
		fmt.Println("Ошибка запуска сервера:", err)
	}
}
