// internal/cache/cache.go

package cache

import (
	"context"
	"encoding/json"
	"os"
	"time"

	"github.com/devonphone/weather-aggregator/internal/models"
	"github.com/redis/go-redis/v9"
)

type Cache interface {
    Get(ctx context.Context, key string) (*models.WeatherData, error)
    Set(ctx context.Context, key string, value *models.WeatherData) error
}

type RedisCache struct {
    client        *redis.Client
    cacheDuration time.Duration
}

func NewRedisCache(cacheDuration time.Duration) (*RedisCache, error) {
    redisAddr := os.Getenv("REDIS_ADDR")
    redisUsername := os.Getenv("REDIS_USERNAME")
    redisPassword := os.Getenv("REDIS_PASSWORD")

    // Using your specific Redis connection details
    client := redis.NewClient(&redis.Options{
        Addr:     redisAddr,
        Username: redisUsername,
        Password: redisPassword,
        DB:       0,
    })

    // Test the connection
    ctx := context.Background()
    _, err := client.Ping(ctx).Result()
    if err != nil {
        return nil, err
    }

    return &RedisCache{
        client:        client,
        cacheDuration: cacheDuration,
    }, nil
}

func (c *RedisCache) Get(ctx context.Context, key string) (*models.WeatherData, error) {
    val, err := c.client.Get(ctx, key).Result()
    if err == redis.Nil {
        return nil, nil
    }
    if err != nil {
        return nil, err
    }

    var weather models.WeatherData
    if err := json.Unmarshal([]byte(val), &weather); err != nil {
        return nil, err
    }

    weather.Cached = true
    return &weather, nil
}

func (c *RedisCache) Set(ctx context.Context, key string, value *models.WeatherData) error {
    data, err := json.Marshal(value)
    if err != nil {
        return err
    }

    return c.client.Set(ctx, key, data, c.cacheDuration).Err()
}

// Add a close method for proper cleanup
func (c *RedisCache) Close() error {
    return c.client.Close()
}