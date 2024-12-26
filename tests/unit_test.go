package tests

import (
	"context"
	"testing"
	"time"

	"github.com/devonphone/weather-aggregator/internal/models"
	"github.com/devonphone/weather-aggregator/internal/providers"
)

func TestOpenWeatherProvider(t *testing.T) {
	provider := providers.NewOpenWeatherProvider("2677f7c214be2d4c307f44a4c6422c1e")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	city := "Jakarta"
	data, err := provider.GetWeather(ctx, city)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if data.City != city {
		t.Errorf("Expected city %s, got %s", city, data.City)
	}

	if data.Temperature == 0 {
		t.Errorf("Expected temperature to be non-zero, got %f", data.Temperature)
	}
}

func TestWeatherAPIProvider(t *testing.T) {
	provider := providers.NewWeatherAPIProvider("56677165c854477d96480002242512")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	city := "Jakarta"
	data, err := provider.GetWeather(ctx, city)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if data.City != city {
		t.Errorf("Expected city %s, got %s", city, data.City)
	}

	if data.Temperature == 0 {
		t.Errorf("Expected temperature to be non-zero, got %f", data.Temperature)
	}
}

func TestFetchFromProviders(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	providersList := []providers.WeatherProvider{
		providers.NewOpenWeatherProvider("2677f7c214be2d4c307f44a4c6422c1e"),
		providers.NewWeatherAPIProvider("56677165c854477d96480002242512"),
	}

	city := "Jakarta"
	results := make(chan interface{}, len(providersList))
	hasError := false

	for _, provider := range providersList {
		go func(p providers.WeatherProvider) {
			data, err := p.GetWeather(ctx, city)
			if err != nil {
				t.Errorf("Error fetching data from provider %s: %v", p.GetProviderName(), err)
				hasError = true
			} else {
				results <- data
			}
		}(provider)
	}

	select {
	case res := <-results:
		weatherData, ok := res.(*models.WeatherData)
		if !ok {
			t.Errorf("Expected *models.WeatherData, got %T", res)
		}
		if weatherData.City != city {
			t.Errorf("Expected city %s, got %s", city, weatherData.City)
		}
	case <-time.After(5 * time.Second):
		if hasError {
			t.Log("Some providers returned errors, but test timed out waiting for valid responses.")
		} else {
			t.Error("Test timed out waiting for provider results")
		}
	}
}
