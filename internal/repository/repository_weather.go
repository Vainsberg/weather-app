package repositoryweather

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/Vainsberg/weather-app/internal/response"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {

	return &Repository{db: db}

}

func (r *Repository) GetWeatherDataByCoordinates(latitudeText, longitudeText string) response.Responseweather {
	row := r.db.QueryRow("SELECT * FROM weatherdata WHERE latitude = $1 AND longitude = $2 AND date >= datetime('now','-1 hours');", latitudeText, longitudeText)

	responseN := response.Responseweather{}

	err := row.Scan(&responseN.Id, &responseN.Date, &responseN.Latitude, &responseN.Longitude, &responseN.Current.Wind, &responseN.Current.Temperature)

	if err != nil && err != sql.ErrNoRows {
		log.Fatal(err)

	}

	return responseN

}

func (r *Repository) GetWeatherAddendumByCoordinates(latitudeText string, longitudeText string, Temperature float64, Wind float64) {
	responseN := response.Responseweather{}
	_, err := r.db.Exec("INSERT INTO weatherdata (date, latitude, longitude, temperature, wind) VALUES (datetime('now'), ?, ?, ?, ?)", latitudeText, longitudeText, responseN.Current.Temperature, responseN.Current.Wind)
	if err != nil {
		log.Fatal(err)
		fmt.Println(err)
	}
}
