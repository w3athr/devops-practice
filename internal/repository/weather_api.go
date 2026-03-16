package repository

import (
	"encoding/json"
	"fmt"
	"net/http"

	"app/internal/entity"
)

type WeatherRepo struct {
	APIKey string
}

func NewWeatherRepo(apiKey string) *WeatherRepo {
	return &WeatherRepo{APIKey: apiKey}
}

// GetRawWeatherData is used for retrieve data via API
func (r *WeatherRepo) GetRawWeatherData(city, from, to string) ([]entity.DayData, error) {
	url := fmt.Sprintf("https://weather.visualcrossing.com/VisualCrossingWebServices/rest/services/timeline/%s", city)
	if from != "" {
		url += "/" + from
		if to != "" {
			url += "/" + to
		}
	}
	url += fmt.Sprintf("?key=%s&unitGroup=metric&include=days&elements=temp,tempmax,tempmin", r.APIKey)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("weather api error: status %d", resp.StatusCode)
	}

	var apiResp struct {
		Days []struct {
			Temp    float64 `json:"temp"`
			TempMax float64 `json:"tempmax"`
			TempMin float64 `json:"tempmin"`
		} `json:"days"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, err
	}

	var result []entity.DayData
	for _, d := range apiResp.Days {
		result = append(result, entity.DayData{
			Temp:    d.Temp,
			TempMin: d.TempMin,
			TempMax: d.TempMax,
		})
	}
	return result, nil
}
