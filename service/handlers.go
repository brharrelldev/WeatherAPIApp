package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/brharrelldev/weatherAPI/constants"
	"github.com/brharrelldev/weatherAPI/models"
	"github.com/brharrelldev/weatherAPI/validations"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"net/http"
	"net/url"
	"strconv"
	"unicode"
)

func (srv *Service) weatherHandler(apiToken string) func(w http.ResponseWriter, r *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {
		var req *models.WeatherRequest

		reqId, err := uuid.NewUUID()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		srv.Logger.Info("starting new request", zap.String("request_id", reqId.String()))
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			srv.Logger.Error("error decoding request", zap.String("request_id", reqId.String()), zap.Error(err))

			se := constants.ServiceError{
				RequestId:  reqId.String(),
				StatusCode: http.StatusBadRequest,
				Err:        err.Error(),
			}

			http.Error(w, se.Error(), http.StatusBadRequest)
			return
		}

		if !unicode.IsUpper(rune(req.State[0])) {
			req.State = cases.Title(language.English, cases.NoLower).String(req.State)
		}

		// validating request. basic city and state validations since this API is only supporting USA
		_, err = validations.ValidateRequest(*req)
		if err != nil {

			// capturing validation errors, a bit berbose
			if errors.Is(err, constants.ErrStateInvalid) {
				srv.Logger.Error("invalid state entered for request", zap.String("request_id", reqId.String()), zap.Error(err))
				se := constants.ServiceError{
					RequestId:  reqId.String(),
					StatusCode: http.StatusBadRequest,
					Err:        constants.ErrStateInvalid.Error(),
				}

				http.Error(w, se.Error(), http.StatusBadRequest)
			}

			if errors.Is(err, constants.ErrCityInvalid) {
				srv.Logger.Error("invalid city entered for request", zap.String("request_id", reqId.String()), zap.Error(err))
				se := constants.ServiceError{
					RequestId:  reqId.String(),
					StatusCode: http.StatusBadRequest,
					Err:        constants.ErrCityInvalid.Error(),
				}

				http.Error(w, se.Error(), http.StatusBadRequest)

			}

			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// building GEO request
		geoRequest, err := url.Parse(srv.WeatherService.GeoLocationEndpoint)
		if err != nil {
			srv.Logger.Error("error parsing request", zap.String("request_id", reqId.String()), zap.Error(err))
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// adding query params
		geoQuery := geoRequest.Query()
		geoQuery.Add("q", fmt.Sprintf("%s,%s,US", req.City, req.State))
		geoQuery.Add("appid", apiToken)
		geoRequest.RawQuery = geoQuery.Encode()

		fmt.Println(geoRequest.String())

		gReq, err := http.NewRequest(http.MethodGet, geoRequest.String(), nil)
		if err != nil {
			srv.Logger.Error("error creating request", zap.String("request_id", reqId.String()), zap.Error(err))
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		client := &http.Client{}
		gResp, err := client.Do(gReq)
		if err != nil {
			srv.Logger.Error("Error getting geolocation data",
				zap.String("request_id", reqId.String()),
				zap.String("status_code", http.StatusText(gResp.StatusCode)),
				zap.Error(err))
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		defer gResp.Body.Close()

		if gResp.StatusCode != http.StatusOK {
			srv.Logger.Error("non-200 response found", zap.String("request_id", reqId.String()), zap.String("status_code", http.StatusText(gResp.StatusCode)))
			http.Error(w, http.StatusText(gResp.StatusCode), gResp.StatusCode)
			return
		}

		srv.Logger.Info("geo request was successful", zap.String("request_id", reqId.String()))

		var geoResp []models.GeoResponse

		if err := json.NewDecoder(gResp.Body).Decode(&geoResp); err != nil {
			srv.Logger.Error("error decoding response", zap.String("request_id", reqId.String()), zap.Error(err))
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		currentWeatherURL, err := url.Parse(srv.WeatherService.CurrentWeatherEndpoint)
		if err != nil {
			srv.Logger.Error("error parsing request", zap.String("request_id", reqId.String()), zap.Error(err))
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		currentQuery := currentWeatherURL.Query()
		currentQuery.Add("lon", strconv.FormatFloat(geoResp[0].Longitude, 'f', -1, 64))
		currentQuery.Add("lat", strconv.FormatFloat(geoResp[0].Latitude, 'f', -1, 64))
		currentQuery.Add("units", "imperial")
		currentQuery.Add("appid", apiToken)
		currentWeatherURL.RawQuery = currentQuery.Encode()

		currentWeatherReq, err := http.NewRequest(http.MethodGet, currentWeatherURL.String(), nil)
		if err != nil {
			srv.Logger.Error("error creating request", zap.String("request_id", reqId.String()), zap.Error(err))
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		wr, err := client.Do(currentWeatherReq)
		if err != nil {
			srv.Logger.Error("Error getting current weather data", zap.String("request_id", reqId.String()), zap.Error(err), zap.Int("status_code", wr.StatusCode))
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if wr.StatusCode != http.StatusOK {
			srv.Logger.Error("non-200 error code found when getting current weather data", zap.String("request_id", reqId.String()), zap.Int("status_code", wr.StatusCode))
			http.Error(w, http.StatusText(wr.StatusCode), wr.StatusCode)
			return
		}

		var currentWeatherResp models.CurrentWeatherResponse

		if err := json.NewDecoder(wr.Body).Decode(&currentWeatherResp); err != nil {
			srv.Logger.Error("error decoding response", zap.String("request_id", reqId.String()), zap.Error(err))
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		defer wr.Body.Close()

		var feelsLike string

		if currentWeatherResp.Main.Temp < 50 {
			feelsLike = "BURRRR!"
		}

		if currentWeatherResp.Main.Temp > 60 {
			feelsLike = "Warm"
		}

		if currentWeatherResp.Main.Temp > 80 {
			feelsLike = "HOT!"
		}

		weatherAPIResponse := models.WeatherResponse{
			Feeling:     feelsLike,
			Temperature: currentWeatherResp.Main.Temp,
		}

		respToUser, err := json.Marshal(weatherAPIResponse)
		if err != nil {
			srv.Logger.Error("error encoding response", zap.String("request_id", reqId.String()), zap.Error(err))
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if _, err := w.Write(respToUser); err != nil {
			srv.Logger.Error("error writing response", zap.String("request_id", reqId.String()), zap.Error(err))
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

}
