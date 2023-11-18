package viper

import (
	"fmt"

	"github.com/spf13/viper"
)

func ValApiKey() string {
	viper.SetConfigFile("config.yaml")

	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println(err)
		return "Ошибка чтения конфигурации:"
	}
	value := viper.GetString("API-KEY")
	return value

}
