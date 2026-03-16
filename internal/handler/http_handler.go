package handler

import (
	"net/http"

	"app/internal/entity"
	"app/internal/usecase"

	"github.com/gin-gonic/gin"
)

type WeatherHandler struct {
	useCase *usecase.WeatherUseCase
	info    entity.ServiceInfo
}

type FinalResponse struct {
	Service string      `json:"service"`
	Data    interface{} `json:"data"`
}

func NewWeatherHandler(uc *usecase.WeatherUseCase, info entity.ServiceInfo) *WeatherHandler {
	return &WeatherHandler{useCase: uc, info: info}
}

func (h *WeatherHandler) GetWeather(c *gin.Context) {
	city := c.Query("city")
	from := c.Query("date_from")
	to := c.Query("date_to")

	if city == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "city is required"})
		return
	}

	stats, err := h.useCase.GetStats(city, from, to)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	finalResponse := FinalResponse{
		Service: h.info.Service, // Service first
		Data: gin.H{
			"temperature_c": stats, // Then data
		},
	}

	c.JSON(http.StatusOK, finalResponse)
}

func (h *WeatherHandler) GetInfo(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, h.info)
}
