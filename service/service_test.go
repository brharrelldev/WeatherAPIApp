package service

import (
	"go.uber.org/zap"
	"reflect"
	"testing"
)

func TestNewService(t *testing.T) {

	tests := []struct {
		name     string
		config   *Config
		expected *Service
	}{
		{
			name: "config",
			config: &Config{
				Host:                   "127.0.0.1",
				Port:                   "8080",
				WeatherBaseURL:         "api.taketest.com",
				CurrentWeatherEndpoint: "fakeEndpoint",
				GeoEndpoint:            "fakeGeo",
			},
			expected: &Service{
				Host: "127.0.0.1",
				Port: "8080",
				WeatherService: WeatherService{
					CurrentWeatherEndpoint: "http://api.taketest.com/fakeEndpoint",
					GeoLocationEndpoint:    "http://api.taketest.com/fakeGeo",
				},
				Logger: nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger := zap.NewNop()
			s, err := NewService(tt.config, logger)
			if err != nil {
				t.Fatal(err)
			}

			if reflect.DeepEqual(s, tt.expected) {
				t.Errorf("got %v, want %v", s, tt.expected)
			}
		})
	}
}
