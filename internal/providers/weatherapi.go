package providers

import (
    "context"
    "encoding/json"
    "fmt"
    "github.com/devonphone/weather-aggregator/internal/models"
    "net/http"
    "time"
)

type WeatherAPIProvider struct {
    apiKey string
    client *http.Client
}

func NewWeatherAPIProvider(apiKey string) *WeatherAPIProvider {
    return &WeatherAPIProvider{
        apiKey: apiKey,
        client: &http.Client{Timeout: 10 * time.Second},
    }
}

func (p *WeatherAPIProvider) GetWeather(ctx context.Context, city string) (*models.WeatherData, error) {
	url := fmt.Sprintf(
		"http://api.weatherapi.com/v1/current.json?key=%s&q=%s&aqi=no", 
		p.apiKey, city,
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
        return nil, fmt.Errorf("WeatherAPI API error: %d", resp.StatusCode)
    }

	var result struct {
        Location struct {
            Name    string `json:"name"`
            Region  string `json:"region"`
            Country string `json:"country"`
        } `json:"location"`
        Current struct {
            TempC      float64 `json:"temp_c"`
            Humidity   int     `json:"humidity"`
            Condition  struct {
                Text string `json:"text"`
            } `json:"condition"`
        } `json:"current"`
    }

    if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
        return nil, err
    }

    return &models.WeatherData{
        City:        result.Location.Name,
        Temperature: result.Current.TempC,
        Humidity:    result.Current.Humidity,
        Condition:   result.Current.Condition.Text,
        Source:      p.GetProviderName(),
        Timestamp:   time.Now(),
    }, nil
}

func (p *WeatherAPIProvider) GetProviderName() string {
    return "WeatherAPIMap"
}