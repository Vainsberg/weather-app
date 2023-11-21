package main

import (
	"fmt"
	"net/http"
	"weather/internal/getweather"
	"weather/internal/viper"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	viper.ValApiKey()
	http.HandleFunc("/get_weather", getweather.Handler)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Ошибка запуска сервера:", err)
	}
}
