package entity

type TemperatureStats struct {
	Average float64 `json:"average"`
	Median  float64 `json:"median"`
	Min     float64 `json:"min"`
	Max     float64 `json:"max"`
}

type DayData struct {
	Temp    float64 `json:"temp"`
	TempMin float64 `json:"tempmin"`
	TempMax float64 `json:"tempmax"`
}
