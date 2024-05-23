package service

import (
	"errors"
	"fmt"
	"github.com/brharrelldev/weatherAPI/constants"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"net"
	"net/http"
	"time"
)

func NewService(c *Config, logger *zap.Logger) (*Service, error) {

	if c.WeatherBaseURL == "" {
		return nil, constants.ErrNoValidWeatherEndpoint
	}

	if c.CurrentWeatherEndpoint == "" {
		return nil, constants.ErrNoValidCurrentWeatherEndpoint
	}

	if c.GeoEndpoint == "" {
		return nil, constants.ErrNoValidGeoEndpoint
	}

	weatherService := WeatherService{
		CurrentWeatherEndpoint: fmt.Sprintf("http://%s/%s", c.WeatherBaseURL, c.CurrentWeatherEndpoint),
		GeoLocationEndpoint:    fmt.Sprintf("http://%s/%s", c.WeatherBaseURL, c.GeoEndpoint),
	}

	if logger == nil {
		return nil, errors.New("logger is null, and we won't start service without it")
	}

	return &Service{
		Host:           c.Host,
		Port:           c.Port,
		WeatherService: weatherService,
		Logger:         logger,
	}, nil

}

func (srv *Service) Start(apiKey string) error {
	host := net.JoinHostPort(srv.Host, srv.Port)

	srv.Logger.Info("server is starting", zap.String("host", host))

	router := mux.NewRouter()

	router.HandleFunc("/weather", srv.weatherHandler(apiKey)).Methods(http.MethodPost)

	service := http.Server{
		Addr:         host,
		Handler:      router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}

	return service.ListenAndServe()
}
