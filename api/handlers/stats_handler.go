package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/devonphone/weather-aggregator/internal/models"
	"github.com/devonphone/weather-aggregator/internal/stats"
)

type StatsHandler struct {
    stats *stats.StatsTracker
}

func NewStatsHandler(stats *stats.StatsTracker) *StatsHandler {
    return &StatsHandler{
        stats: stats,
    }
}

func (h *StatsHandler) GetStats(w http.ResponseWriter, r *http.Request) {
    stats := h.stats.GetStats()
    RespondJSON(w, stats)
}

// Common response utilities
func RespondJSON(w http.ResponseWriter, data interface{}) {
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(data)
}

func RespondError(w http.ResponseWriter, code int, message string) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(code)
    json.NewEncoder(w).Encode(models.ErrorResponse{
        Error:   http.StatusText(code),
        Code:    code,
        Message: message,
    })
}
