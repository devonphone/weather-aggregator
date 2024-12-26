package config

import (
    "github.com/joho/godotenv"
    "os"
    "strconv"
    "time"
)

type Config struct {
    Port                string
    OpenWeatherApiKey   string
    WeatherApiKey       string
    RedisUrl           string
    CacheDuration      time.Duration
    RateLimitRequests  int
    RateLimitDuration  time.Duration
}

func LoadConfig() (*Config, error) {
    if err := godotenv.Load(); err != nil {
        return nil, err
    }

    rateLimitReq, _ := strconv.Atoi(getEnv("RATE_LIMIT_REQUESTS", "60"))
    cacheDuration, _ := time.ParseDuration(getEnv("CACHE_DURATION", "30m"))
    rateLimitDuration, _ := time.ParseDuration(getEnv("RATE_LIMIT_DURATION", "1m"))

    return &Config{
        Port:               getEnv("PORT", "8080"),
        OpenWeatherApiKey:  os.Getenv("OPENWEATHER_API_KEY"),
        WeatherApiKey:      os.Getenv("WEATHERAPI_KEY"),
        RedisUrl:          getEnv("REDIS_URL", "redis://localhost:6379"),
        CacheDuration:     cacheDuration,
        RateLimitRequests: rateLimitReq,
        RateLimitDuration: rateLimitDuration,
    }, nil
}

func getEnv(key, defaultValue string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return defaultValue
}