package infra

import (
	"time"

	"github.com/go-redis/redis/v8"
)

// RedisConfig holds Redis connection configuration
type RedisConfig struct {
	Host     string
	Port     int
	Password string
	DB       int
}

// NewRedisClient creates a new Redis client for idempotency
func NewRedisClient(config RedisConfig) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     config.Host + ":" + string(rune(config.Port)),
		Password: config.Password,
		DB:       config.DB,
	})
}

// IdempotencyKey represents an idempotent operation
type IdempotencyKey struct {
	Key       string
	Value     string
	ExpiresAt time.Duration
}
