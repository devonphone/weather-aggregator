package providers

import (
    "context"
    "encoding/json"
    "fmt"
    "github.com/devonphone/weather-aggregator/internal/models"
    "net/http"
    "time"
)

type OpenWeatherProvider struct {
    apiKey string
    client *http.Client
}

func NewOpenWeatherProvider(apiKey string) *OpenWeatherProvider {
    return &OpenWeatherProvider{
        apiKey: apiKey,
        client: &http.Client{Timeout: 10 * time.Second},
    }
}

func (p *OpenWeatherProvider) GetWeather(ctx context.Context, city string) (*models.WeatherData, error) {
    url := fmt.Sprintf(
        "https://api.openweathermap.org/data/2.5/weather?q=%s&appid=%s&units=metric",
        city, p.apiKey,
    )

    req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
    if err != nil {
        return nil, err
    }

    resp, err := p.client.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("OpenWeather API error: %d", resp.StatusCode)
    }

    var result struct {
        Main struct {
            Temp     float64 `json:"temp"`
            Humidity int     `json:"humidity"`
        } `json:"main"`
        Weather []struct {
            Description string `json:"description"`
        } `json:"weather"`
    }

    if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
        return nil, err
    }

    return &models.WeatherData{
        City:        city,
        Temperature: result.Main.Temp,
        Humidity:    result.Main.Humidity,
        Condition:   result.Weather[0].Description,
        Source:      p.GetProviderName(),
        Timestamp:   time.Now(),
    }, nil
}

func (p *OpenWeatherProvider) GetProviderName() string {
    return "OpenWeatherMap"
}