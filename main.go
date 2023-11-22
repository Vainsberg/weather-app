package main

import (
	"fmt"
	"net/http"

	repositoryweather "github.com/Vainsberg/weather-app/internal/repository"
	_ "github.com/mattn/go-sqlite3"
)

func main() {

	db := db.NewDB()
	defer db.Close()

	repository := repositoryweather.NewRepository(db)
	handler := repositoryweather.Repository(repository)

	http.HandleFunc("/get_weather", handler.GetWeather)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Ошибка запуска сервера:", err)
	}
}
