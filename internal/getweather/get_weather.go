package getweather

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	repositoryweather "github.com/Vainsberg/weather-app/internal/repository"
	"github.com/Vainsberg/weather-app/internal/response"
	formatint "github.com/Vainsberg/weather-app/pkg/formatInt"
)

type Handler struct {
	weatherRepository repositoryweather.Repository
}

func NewHandler(repos *repositoryweather.Repository) *Handler {
	return &Handler{weatherRepository: *repos}
}

func (h *Handler) GetWeather(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	latitudeText := query.Get("latitude")
	longitudeText := query.Get("longitude")
	responseN := h.weatherRepository.GetWeatherDataByCoordinates(latitudeText, longitudeText)

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
		var responseWeather response.Responseweather
		err = json.Unmarshal(data, &responseWeather)
		if err != nil {
			log.Println("Error parsing weather data:", err)
			http.Error(w, "Error parsing weather data", http.StatusInternalServerError)
			fmt.Println(err)

			return

		}
		responseText := fmt.Sprintf("Добрый день! Сегодня температура %d градусов, скорость ветра %d м/с.", formatint.FormatInt(responseWeather.Current.Temperature), formatint.FormatInt(responseWeather.Current.Temperature))
		h.weatherRepository.GetWeatherAddendumByCoordinates(latitudeText, longitudeText, responseN.Current.Temperature, responseN.Current.Wind)
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(responseText))
	}
}
