package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

type responseweather struct {
	Id        int
	Date      string
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Current   struct {
		Wind        float64 `json:"wind_speed_10m"`
		Temperature float64 `json:"temperature_2m"`
	} `json:"current"`
}

func handler(w http.ResponseWriter, r *http.Request) {

	query := r.URL.Query()
	latitudeText := query.Get("latitude")
	longitudeText := query.Get("longitude")

	row := db.QueryRow("SELECT * FROM weatherdata WHERE latitude = $1 AND longitude = $2 AND date >= datetime('now','-1 hours');", latitudeText, longitudeText)

	responseN := responseweather{}

	err := row.Scan(&responseN.Id, &responseN.Date, &responseN.Latitude, &responseN.Longitude, &responseN.Current.Wind, &responseN.Current.Temperature)

	if err != nil && err != sql.ErrNoRows {
		log.Fatal(err)

	}

	if responseN.Latitude != 0 || responseN.Longitude != 0 {
		responseText := fmt.Sprintf("Добрый день! Сегодня температура %d градусов, скорость ветра %d м/с.", formatInt(responseN.Current.Temperature), formatInt(responseN.Current.Wind))

		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(responseText))
	} else {

		weatherAPIURL := "https://api.open-meteo.com/v1/forecast?latitude=" + latitudeText + "&longitude=" + longitudeText + "&current=temperature_2m,wind_speed_10m"
		weatherResp, err := http.Get(weatherAPIURL)
		if err != nil {
			log.Println("Error fetching weather data:", err)
			http.Error(w, "Error fetching weather data", http.StatusInternalServerError)
			fmt.Println(err)
			return
		}
		defer weatherResp.Body.Close()

		data, err := io.ReadAll(weatherResp.Body)
		if err != nil {
			log.Println(err)
			fmt.Println(err)
			return
		}
		var responseWeather responseweather
		err = json.Unmarshal(data, &responseWeather)
		if err != nil {
			log.Println("Error parsing weather data:", err)
			http.Error(w, "Error parsing weather data", http.StatusInternalServerError)
			fmt.Println(err)
			return
		}
		responseText := fmt.Sprintf("Добрый день! Сегодня температура %d градусов, скорость ветра %d м/с.", formatInt(responseWeather.Current.Temperature), formatInt(responseWeather.Current.Temperature))
		_, err = db.Exec("INSERT INTO weatherdata (date, latitude, longitude, temperature, wind) VALUES (datetime('now'), ?, ?, ?, ?)", latitudeText, longitudeText, responseWeather.Current.Temperature, responseWeather.Current.Wind)

		if err != nil {
			log.Fatal(err)
			fmt.Println(err)
		}
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(responseText))

	}

}

func formatInt(num float64) int {
	return int(math.Ceil(num))
}

func main() {
	var err error

	db, err = sql.Open("sqlite3", "store.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS weatherdata (
			id INTEGER PRIMARY KEY,
		    date TEXT,
			latitude  FLOAT,
			longitude FLOAT,
			temperature FLOAT,
			wind FLOAT
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
