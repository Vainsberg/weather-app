package getweather

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	formatint "weather/formatInt"
)

var db *sql.DB

func Handler(w http.ResponseWriter, r *http.Request) {
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
		responseText := fmt.Sprintf("Добрый день! Сегодня температура %d градусов, скорость ветра %d м/с.", formatint.FormatInt(responseN.Current.Temperature), formatint.FormatInt(responseN.Current.Wind))

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
		responseText := fmt.Sprintf("Добрый день! Сегодня температура %d градусов, скорость ветра %d м/с.", formatint.FormatInt(responseWeather.Current.Temperature), formatint.FormatInt(responseWeather.Current.Temperature))
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
