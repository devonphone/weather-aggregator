package providers

import (
    "context"
    "github.com/devonphone/weather-aggregator/internal/models"
)

type WeatherProvider interface {
    GetWeather(ctx context.Context, city string) (*models.WeatherData, error)
    GetProviderName() string
}