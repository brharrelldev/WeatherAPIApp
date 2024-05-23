package models

import "go.uber.org/zap"

type Service struct {
	Host           string
	Port           string
	WeatherService WeatherService
	Logger         *zap.Logger
}

type WeatherService struct {
	CurrentWeatherEndpoint string
	GeoLocationEndpoint    string
}

type GeoResponse struct {
	Name       string            `json:"name,omitempty"`
	LocalNames map[string]string `json:"local_names,omitempty"`
	Latitude   float64           `json:"latitude,omitempty"`
	Longitude  float64           `json:"longitude,omitempty"`
	Country    string            `json:"country,omitempty"`
	State      string            `json:"state,omitempty"`
}

type GeoResponses []GeoResponse

type WeatherResponse struct {
	Feeling     string  `json:"feeling,omitempty"`
	Temperature float64 `json:"temperature,omitempty"`
}

type CurrentWeatherResponse struct {
	Coord      *Coord    `json:"coord,omitempty"`
	Weather    []Weather `json:"weather,omitempty"`
	Main       *Main     `json:"main,omitempty"`
	Wind       *Wind     `json:"wind,omitempty"`
	Clouds     *Clouds   `json:"clouds,omitempty"`
	Sys        *Sys      `json:"sys,omitempty"`
	Dt         int32     `json:"dt,omitempty"`
	Visibility int       `json:"visibility,omitempty"`
	Id         int32     `json:"id,omitempty"`
	Timezone   int32     `json:"timezone,omitempty"`
	Base       string    `json:"base,omitempty"`
	Name       string    `json:"name"`
	Cod        int       `json:"cod"`
}

type Coord struct {
	Lon float64 `json:"lon,omitempty"`
	Lat float64 `json:"lat,omitempty"`
}

type Weather struct {
	Id          int    `json:"id,omitempty"`
	Main        string `json:"main,omitempty"`
	Description string `json:"description,omitempty"`
	Icon        string `json:"icon,omitempty"`
}

type Main struct {
	Temp      float64 `json:"temp,omitempty"`
	FeelsLike float64 `json:"feels_like,omitempty"`
	TempMin   float64 `json:"temp_min,omitempty"`
	TempMax   float64 `json:"temp_max,omitempty"`
	Pressure  int     `json:"pressure,omitempty"`
	Humidity  int     `json:"humidity,omitempty"`
}

type Wind struct {
	Speed float64 `json:"speed,omitempty"`
	Deg   int     `json:"deg,omitempty"`
}

type Clouds struct {
	All int `json:"all,omitempty"`
}

type Sys struct {
	Type    int    `json:"type,omitempty"`
	ID      int    `json:"id,omitempty"`
	Country string `json:"country,omitempty"`
	Sunrise int32  `json:"sunrise,omitempty"`
	Sunset  int32  `json:"sunset,omitempty"`
}
