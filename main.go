// main.go

package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/devonphone/weather-aggregator/api/handlers"
	"github.com/devonphone/weather-aggregator/api/routes"
	"github.com/devonphone/weather-aggregator/config"
	"github.com/devonphone/weather-aggregator/internal/cache"
	"github.com/devonphone/weather-aggregator/internal/providers"
	"github.com/devonphone/weather-aggregator/internal/rate_limit"
	"github.com/devonphone/weather-aggregator/internal/stats"
	"github.com/gorilla/mux"
)

func main() {
    // Load configuration
    cfg, err := config.LoadConfig()
    if err != nil {
        log.Fatalf("Failed to load config: %v", err)
    }

    // Initialize cache
    weatherCache, err := cache.NewRedisCache(cfg.CacheDuration)
    if err != nil {
        log.Fatalf("Failed to initialize cache: %v", err)
    }
    defer weatherCache.Close()

    // Initialize weather providers
    weatherProviders := []providers.WeatherProvider{
		providers.NewWeatherAPIProvider(cfg.WeatherApiKey),
        providers.NewOpenWeatherProvider(cfg.OpenWeatherApiKey),
    }

    // Initialize rate limiter
    rateLimiter := ratelimit.NewTokenBucketLimiter(
        cfg.RateLimitRequests,
        cfg.RateLimitDuration,
    )

    // Initialize stats tracker
    statsTracker := stats.NewStatsTracker()

    // Initialize handlers
    weatherHandler := handlers.NewWeatherHandler(
        weatherProviders,
        weatherCache,
        rateLimiter,
        statsTracker,
    )
    statsHandler := handlers.NewStatsHandler(statsTracker)

    // Initialize router and API
    router := mux.NewRouter()
    api := routes.NewAPI(weatherHandler, statsHandler)
    api.SetupRoutes(router)

    // Create server
    server := &http.Server{
        Addr:         ":" + cfg.Port,
        Handler:      router,
        ReadTimeout:  15 * time.Second,
        WriteTimeout: 15 * time.Second,
        IdleTimeout:  60 * time.Second,
    }

    // Channel to listen for errors coming from the server
    serverErrors := make(chan error, 1)

    // Start server
    go func() {
        log.Printf("Server is starting on port %s...\n", cfg.Port)
        serverErrors <- server.ListenAndServe()
    }()

    // Channel to listen for interrupt signals
    shutdown := make(chan os.Signal, 1)
    signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

    // Block until we receive a signal or error
    select {
    case err := <-serverErrors:
        log.Fatalf("Error starting server: %v", err)

    case <-shutdown:
        log.Println("Starting shutdown...")
        
        // Create shutdown context with timeout
        ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
        defer cancel()

        // Shutdown server gracefully
        err := server.Shutdown(ctx)
        if err != nil {
            log.Printf("Error during server shutdown: %v", err)
            err = server.Close()
        }

        switch {
        case err != nil:
            log.Fatalf("Error during server shutdown: %v", err)
        case ctx.Err() != nil:
            log.Fatalf("Timeout during shutdown: %v", ctx.Err())
        }
    }
}