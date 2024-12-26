package models

import "time"

type WeatherData struct {
    City        string    `json:"city"`
    Temperature float64   `json:"temperature"`
    Humidity    int       `json:"humidity"`
    Condition   string    `json:"condition"`
    Source      string    `json:"source"`
    Cached      bool      `json:"cached"`
    Timestamp   time.Time `json:"timestamp"`
}

type ErrorResponse struct {
    Error   string `json:"error"`
    Code    int    `json:"code"`
    Message string `json:"message"`
}

type StatsResponse struct {
    TotalRequests  int64 `json:"total_requests"`
    CacheHits      int64 `json:"cache_hits"`
    CacheMisses    int64 `json:"cache_misses"`
    ApiCalls       int64 `json:"api_calls"`
    RateLimitHits  int64 `json:"rate_limit_hits"`
}