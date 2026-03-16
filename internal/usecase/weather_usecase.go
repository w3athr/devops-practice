package usecase

import (
	"math"
	"sort"

	"app/internal/entity"
)

type WeatherRepository interface {
	GetRawWeatherData(city, from, to string) ([]entity.DayData, error)
}

type WeatherUseCase struct {
	repo WeatherRepository
}

func NewWeatherUseCase(repo WeatherRepository) *WeatherUseCase {
	return &WeatherUseCase{repo: repo}
}

func (uc *WeatherUseCase) GetStats(city, from, to string) (entity.TemperatureStats, error) {
	days, err := uc.repo.GetRawWeatherData(city, from, to)
	if err != nil || len(days) == 0 {
		return entity.TemperatureStats{}, err
	}

	var sum float64
	var temps []float64

	minTemp := days[0].TempMin
	maxTemp := days[0].TempMax

	for _, d := range days {
		sum += d.Temp
		temps = append(temps, d.Temp)

		if d.TempMin < minTemp {
			minTemp = d.TempMin
		}
		if d.TempMax > maxTemp {
			maxTemp = d.TempMax
		}
	}

	sort.Float64s(temps)
	n := len(temps)
	var median float64
	if n%2 == 0 {
		median = (temps[n/2-1] + temps[n/2]) / 2
	} else {
		median = temps[n/2]
	}

	return entity.TemperatureStats{
		Average: round(sum/float64(n), 2),
		Median:  round(median, 2),
		Min:     minTemp,
		Max:     maxTemp,
	}, nil
}

func round(val float64, precision int) float64 {
	ratio := math.Pow(10, float64(precision))
	return math.Round(val*ratio) / ratio
}
