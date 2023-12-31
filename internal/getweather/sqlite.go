package getweather

import (
	"database/sql"
	"log"
)

var db *sql.DB

func CreateDB() *sql.DB {
	var err error
	db, err = sql.Open("sqlite3", "store.db")
	if err != nil {
		log.Fatal(err)
	}

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

	return db
}
