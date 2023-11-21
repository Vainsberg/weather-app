package getweather

type Responseweather struct {
	Id        int
	Date      string
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Current   struct {
		Wind        float64 `json:"wind_speed_10m"`
		Temperature float64 `json:"temperature_2m"`
	} `json:"current"`
}
