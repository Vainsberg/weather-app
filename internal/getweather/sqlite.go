package getweather

import (
	"database/sql"
	"log"
)

func sqlite() {
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
}
