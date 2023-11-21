package repositoryweather

import (
	"database/sql"
	"log"
	"weather/internal/getweather"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) GetWeatherDataByCoordinates(latitudeText, longitudeText string) getweather.Responseweather {
	row := r.db.QueryRow("SELECT * FROM weatherdata WHERE latitude = $1 AND longitude = $2 AND date >= datetime('now','-1 hours');", latitudeText, longitudeText)

	responseN := getweather.Responseweather{}

	err := row.Scan(&responseN.Id, &responseN.Date, &responseN.Latitude, &responseN.Longitude, &responseN.Current.Wind, &responseN.Current.Temperature)

	if err != nil && err != sql.ErrNoRows {
		log.Fatal(err)

	}
	return responseN
}
