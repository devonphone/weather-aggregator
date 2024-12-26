package stats

import (
    "sync/atomic"
    "github.com/devonphone/weather-aggregator/internal/models"
)

type StatsTracker struct {
    totalRequests  int64
    cacheHits      int64
    cacheMisses    int64
    apiCalls       int64
    rateLimitHits  int64
}

func NewStatsTracker() *StatsTracker {
    return &StatsTracker{}
}

func (s *StatsTracker) IncrementRequests() {
    atomic.AddInt64(&s.totalRequests, 1)
}

func (s *StatsTracker) IncrementCacheHits() {
    atomic.AddInt64(&s.cacheHits, 1)
}

func (s *StatsTracker) IncrementCacheMisses() {
    atomic.AddInt64(&s.cacheMisses, 1)
}

func (s *StatsTracker) IncrementApiCalls() {
    atomic.AddInt64(&s.apiCalls, 1)
}

func (s *StatsTracker) IncrementRateLimitHits() {
    atomic.AddInt64(&s.rateLimitHits, 1)
}

func (s *StatsTracker) GetStats() models.StatsResponse {
    return models.StatsResponse{
        TotalRequests:  atomic.LoadInt64(&s.totalRequests),
        CacheHits:      atomic.LoadInt64(&s.cacheHits),
        CacheMisses:    atomic.LoadInt64(&s.cacheMisses),
        ApiCalls:       atomic.LoadInt64(&s.apiCalls),
        RateLimitHits:  atomic.LoadInt64(&s.rateLimitHits),
    }
}