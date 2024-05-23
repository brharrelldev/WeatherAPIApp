package models

type WeatherRequest struct {
	City  string `json:"city"`
	State string `json:"state"`
}
