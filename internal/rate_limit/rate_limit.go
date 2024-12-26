package ratelimit

import (
    "context"
    "golang.org/x/time/rate"

    "time"
)

type RateLimiter interface {
    Allow() bool
    Wait(context.Context) error
}

type TokenBucketLimiter struct {
    limiter *rate.Limiter
}

func NewTokenBucketLimiter(requests int, duration time.Duration) *TokenBucketLimiter {
    return &TokenBucketLimiter{
        limiter: rate.NewLimiter(rate.Every(duration/time.Duration(requests)), requests),
    }
}

func (t *TokenBucketLimiter) Allow() bool {
    return t.limiter.Allow()
}

func (t *TokenBucketLimiter) Wait(ctx context.Context) error {
    return t.limiter.Wait(ctx)
}