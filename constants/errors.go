package constants

import (
	"encoding/json"
	"errors"
)

var (
	ErrNoValidCurrentWeatherEndpoint error = errors.New("no valid current weather endpoint")
	ErrNoValidGeoEndpoint            error = errors.New("no valid geo endpoint")
	ErrNoValidWeatherEndpoint        error = errors.New("no valid weather endpoint")
	ErrStateInvalid                  error = errors.New("invalid state entered")
	ErrCityInvalid                   error = errors.New("invalid city entered")
)
var _ error = &ServiceError{}

type ServiceError struct {
	RequestId  string
	StatusCode int
	Err        string
}

func (s ServiceError) Error() string {
	//TODO implement me
	se, _ := json.Marshal(s)

	return string(se)
}
