package service

import "go.uber.org/zap"

type Config struct {
	Host                   string `yaml:"host"`
	Port                   string `yaml:"port"`
	WeatherBaseURL         string `yaml:"weather_base_url"`
	CurrentWeatherEndpoint string `yaml:"current_weather_endpoint"`
	GeoEndpoint            string `yaml:"geo_endpoint"`
}

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
