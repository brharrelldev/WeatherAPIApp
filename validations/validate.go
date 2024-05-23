package validations

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"fmt"
	"github.com/brharrelldev/weatherAPI/constants"
	"github.com/brharrelldev/weatherAPI/models"
	"strings"
)

//go:embed city_data/cities.json
var cities []byte

func ValidateRequest(req models.WeatherRequest) (bool, error) {

	cityMap := make(map[string][]string)

	if err := json.NewDecoder(bytes.NewBuffer(cities)).Decode(&cityMap); err != nil {
		return false, fmt.Errorf("error decoding city map: %v", err)
	}

	var cityMatch bool
	if val, ok := cityMap[req.State]; ok {

		for _, city := range val {
			if strings.ToLower(city) == strings.ToLower(req.City) {
				cityMatch = true
				break
			}
		}
	} else {
		return false, constants.ErrStateInvalid
	}

	if !cityMatch {
		return false, constants.ErrCityInvalid
	}

	return cityMatch, nil

}
