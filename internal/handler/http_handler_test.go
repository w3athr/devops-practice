package handler

import (
	"app/internal/entity"
	"app/internal/usecase"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

type MockWeatherRepo struct{}

func (m *MockWeatherRepo) GetRawWeatherData(city, from, to string) ([]entity.DayData, error) {
	switch city {
	case "Arkhangelsk":
		return []entity.DayData{
			{Temp: -8.8, TempMin: -8.2, TempMax: -4},
			{Temp: -4.0, TempMin: -8.2, TempMax: -4},
		}, nil
	case "MOSKVA":
		return []entity.DayData{
			{Temp: 4.4, TempMin: 3.8, TempMax: 13.5},
			{Temp: 13.5, TempMin: 3.8, TempMax: 13.5},
		}, nil
	case "Saint-Petersburg":
		return []entity.DayData{
			{Temp: -3.5, TempMin: -1.2, TempMax: 8},
			{Temp: 8.0, TempMin: -1.2, TempMax: 8},
		}, nil
	}
	return nil, fmt.Errorf("city not found in mock")
}

func TestWeatherApp_E2E(t *testing.T) {
	gin.SetMode(gin.TestMode)

	serviceInfo := entity.ServiceInfo{
		Version: "0.5.3",
		Service: "weather",
		Author:  "egor.volkov",
	}

	mockRepo := &MockWeatherRepo{}
	uc := usecase.NewWeatherUseCase(mockRepo)
	h := NewWeatherHandler(uc, serviceInfo)

	r := gin.New()
	r.GET("/info", h.GetInfo)
	r.GET("/info/weather", h.GetWeather)

	t.Run("Endpoint /info", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/info", nil)
		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected 200, got %d", w.Code)
		}

		var res entity.ServiceInfo
		json.Unmarshal(w.Body.Bytes(), &res)

		if res.Service != "weather" || res.Author != "egor.volkov" {
			t.Errorf("Unexpected info output: %s", w.Body.String())
		}
	})

	weatherTests := []struct {
		url      string
		expected string
	}{
		{
			url:      "/info/weather?city=Arkhangelsk&date_from=2021-01-01&date_to=2021-01-02",
			expected: `{"service":"weather","data":{"temperature_c":{"average":-6.4,"median":-6.4,"min":-8.2,"max":-4}}}`,
		},
		{
			url:      "/info/weather?city=MOSKVA&date_from=2026-03-13&date_to=2026-03-14",
			expected: `{"service":"weather","data":{"temperature_c":{"average":8.95,"median":8.95,"min":3.8,"max":13.5}}}`,
		},
		{
			url:      "/info/weather?city=Saint-Petersburg&date_from=2024-03-25&date_to=2024-03-26",
			expected: `{"service":"weather","data":{"temperature_c":{"average":2.25,"median":2.25,"min":-1.2,"max":8}}}`,
		},
	}

	for _, tt := range weatherTests {
		t.Run(tt.url, func(t *testing.T) {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", tt.url, nil)
			r.ServeHTTP(w, req)

			if w.Code != http.StatusOK {
				t.Errorf("Expected 200, got %d. Body: %s", w.Code, w.Body.String())
			}

			var got, want interface{}
			json.Unmarshal(w.Body.Bytes(), &got)
			json.Unmarshal([]byte(tt.expected), &want)

			gotStr, _ := json.Marshal(got)
			wantStr, _ := json.Marshal(want)

			if string(gotStr) != string(wantStr) {
				t.Errorf("\nRequest: %s\nGot:  %s\nWant: %s", tt.url, string(gotStr), string(wantStr))
			}
		})
	}
}
