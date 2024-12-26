package handlers

import (
	"context"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/devonphone/weather-aggregator/internal/cache"
	"github.com/devonphone/weather-aggregator/internal/models"
	"github.com/devonphone/weather-aggregator/internal/providers"
	ratelimit "github.com/devonphone/weather-aggregator/internal/rate_limit"
	"github.com/devonphone/weather-aggregator/internal/stats"
)

type WeatherHandler struct {
	providers   []providers.WeatherProvider
	cache       cache.Cache
	rateLimiter ratelimit.RateLimiter
	stats       *stats.StatsTracker
}

func NewWeatherHandler(
	providers []providers.WeatherProvider,
	cache cache.Cache,
	rateLimiter ratelimit.RateLimiter,
	stats *stats.StatsTracker,
) *WeatherHandler {
	return &WeatherHandler{
		providers:   providers,
		cache:       cache,
		rateLimiter: rateLimiter,
		stats:       stats,
	}
}

func (h *WeatherHandler) GetWeather(w http.ResponseWriter, r *http.Request) {
	h.stats.IncrementRequests()

	city := r.URL.Query().Get("city")
	if city == "" {
		RespondError(w, http.StatusBadRequest, "City parameter is required")
		return
	}

	// Check cache first
	if weatherData, err := h.cache.Get(r.Context(), city); err == nil && weatherData != nil {
		h.stats.IncrementCacheHits()
		log.Printf("[DEBUG] Cache hit for city: %s", city)
		RespondJSON(w, weatherData)
		return
	}
	h.stats.IncrementCacheMisses()
	log.Printf("[DEBUG] Cache miss for city: %s", city)

	// Check rate limit
	if !h.rateLimiter.Allow() {
		h.stats.IncrementRateLimitHits()
		RespondError(w, http.StatusTooManyRequests, "Rate limit exceeded")
		return
	}

	// Log when fetching data from providers
	log.Printf("[DEBUG] Fetching weather data from providers for city: %s", city)

	// Fetch from all providers concurrently
	weatherData := h.fetchFromProviders(r.Context(), city)
	if weatherData == nil {
		RespondError(w, http.StatusServiceUnavailable, "Failed to fetch weather data")
		return
	}

	data, ok := weatherData.(*models.WeatherData)
	if !ok {
		RespondError(w, http.StatusInternalServerError, "Unexpected data type returned from fetchFromProviders")
		return
	}

	// Cache the result
	if err := h.cache.Set(r.Context(), city, data); err != nil {
		log.Printf("[ERROR] Cache set error: %v", err)
	}

	log.Printf("[INFO] Weather data fetched successfully for city: %s", city)
	RespondJSON(w, weatherData)
}

func (h *WeatherHandler) fetchFromProviders(ctx context.Context, city string) interface{} {
	var wg sync.WaitGroup
	results := make(chan interface{}, len(h.providers))
	errors := make(chan error, len(h.providers))

	// Fetch data from each provider concurrently
	for _, provider := range h.providers {
		wg.Add(1)
		go func(provider providers.WeatherProvider) {
			defer wg.Done()
			log.Printf("Fetching data from provider: %s for city: %s\n", provider.GetProviderName(), city)

			// Increment API call count
			h.stats.IncrementApiCalls()
			weatherData, err := provider.GetWeather(ctx, city)
			if err != nil {
				log.Printf("Error from provider %s: %v\n", provider.GetProviderName(), err)
				errors <- err
				return
			}
			log.Printf("Received data from provider %s: %+v\n", provider.GetProviderName(), weatherData)
			results <- weatherData
		}(provider)
	}

	// Close the channels when all goroutines finish
	go func() {
		wg.Wait()
		close(results)
		close(errors)
	}()

	// Wait for the first successful result or return nil if all fail
	select {
	case weatherData := <-results:
		log.Println("Successfully fetched data from a provider.")
		return weatherData
	case <-time.After(5 * time.Second):
		log.Println("All providers timed out")
		return nil
	case err := <-errors:
		log.Printf("Provider error occurred: %v\n", err)
		return nil
	}
}
