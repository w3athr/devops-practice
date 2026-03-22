package main

import (
	"log"
	"os"
	"strconv"
	"strings"

	"app/internal/entity"
	"app/internal/handler"
	"app/internal/repository"
	"app/internal/usecase"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()

	// 1. Инициализация конфига/инфо
	info := entity.ServiceInfo{
		Version: getEnv("VERSION", "1.0.0"),
		Service: getEnv("SERVICE", "weather"),
		Author:  getEnv("AUTHOR", "egor.volkov"),
	}

	// 2. Сборка слоев (Dependency Injection)
	repo := repository.NewWeatherRepo(getAPIKey())
	uc := usecase.NewWeatherUseCase(repo)
	h := handler.NewWeatherHandler(uc, info)

	// 3. Запуск сервера
	r := gin.Default()
	r.GET("/info", h.GetInfo)
	r.GET("/info/weather", h.GetWeather)

	port := getEnv("PORT", "8000")

	if _, err := strconv.Atoi(port); err != nil {
		log.Fatalf("ERROR: PORT environment variable must be a number, but got '%s'", port)
	}

	r.Run(":" + port)
}

func getEnv(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}

func getAPIKey() string {
	// 1. check is exists in docker secrets (for gitlab CI/CD)
	secretPath := "/run/secrets/weather_api_key"
	if data, err := os.ReadFile(secretPath); err == nil {
		return strings.TrimSpace(string(data))
	}

	// 2. check if exists in .env
	if key := os.Getenv("API_KEY"); key != "" {
		return key
	}

	// 3. check if exists in file
	if data, err := os.ReadFile("api_key.txt"); err == nil {
		return strings.TrimSpace(string(data))
	}

	return ""
}
