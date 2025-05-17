package database

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/melkdesousa/gamgo/config"
	"github.com/redis/go-redis/v9"
)

var (
	cache     *redis.Client
	cacheOnce sync.Once
)

// GetCacheConnection returns a singleton Redis client instance
func GetCacheConnection() *redis.Client {
	cacheOnce.Do(func() {
		cache = redis.NewClient(&redis.Options{
			Addr:     config.MustGetEnv("CACHE_ADDR"),      // Redis server address
			Password: config.MustGetEnv("CACHE_PASSWORD"),  // No password set
			DB:       config.MustGetEnvAs[int]("CACHE_DB"), // Use default DB
		})

		// Start a goroutine to periodically check cache health
		go startCacheHealthCheck()
	})

	return cache
}

// startCacheHealthCheck runs a periodic health check on the Redis connection
func startCacheHealthCheck() {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		if cache == nil {
			continue
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		_, err := cache.Ping(ctx).Result()
		cancel()

		if err != nil {
			// Log the error
			log.Printf("Cache health check failed: %v", err)

			// Try to reconnect
			reconnectCache()
		}
	}
}

// reconnectCache attempts to reconnect to Redis
func reconnectCache() {
	if cache != nil {
		// Close existing connection
		_ = cache.Close()
	}

	cache = redis.NewClient(&redis.Options{
		Addr:     config.MustGetEnv("CACHE_ADDR"),
		Password: config.MustGetEnv("CACHE_PASSWORD"),
		DB:       config.MustGetEnvAs[int]("CACHE_DB"),
	})

	// Test the new connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if _, err := cache.Ping(ctx).Result(); err != nil {
		log.Printf("Failed to reconnect to cache: %v", err)
	} else {
		log.Println("Successfully reconnected to cache")
	}
}

// CloseCacheConnection closes the Redis client connection
func CloseCacheConnection() {
	if cache != nil {
		err := cache.Close()
		if err != nil {
			// Handle error if needed
		}
		cache = nil
	}
}
