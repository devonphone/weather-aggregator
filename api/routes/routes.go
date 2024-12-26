package routes

import (
	"log"
	"net/http"

	"github.com/devonphone/weather-aggregator/api/handlers"
	"github.com/gorilla/mux"
)

type API struct {
    weatherHandler *handlers.WeatherHandler
    statsHandler  *handlers.StatsHandler
}

func NewAPI(
    weatherHandler *handlers.WeatherHandler,
    statsHandler *handlers.StatsHandler,
) *API {
    return &API{
        weatherHandler: weatherHandler,
        statsHandler:  statsHandler,
    }
}

func (api *API) SetupRoutes(router *mux.Router) {
    // Weather endpoints
    router.HandleFunc("/weather", api.weatherHandler.GetWeather).Methods("GET")
    
    // Stats endpoints
    router.HandleFunc("/stats", api.statsHandler.GetStats).Methods("GET")
    
    // Add middleware for all routes
    router.Use(api.loggingMiddleware)
}

func (api *API) loggingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Log the request
        log.Printf(
            "[%s] %s %s",
            r.Method,
            r.RequestURI,
            r.RemoteAddr,
        )
        next.ServeHTTP(w, r)
    })
}