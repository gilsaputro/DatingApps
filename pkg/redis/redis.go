package redis

import (
	"context"
	"log"
	"time"

	redis "github.com/go-redis/redis/v8"
)

// RedisConfig is list config to create Redis client
type RedisConfig struct {
	Host     string
	Port     string
	Password string
}

// RedisMethod is list all available method for redis
type RedisMethod interface {
	Set(key string, value interface{}, expiration time.Duration) error
	Get(key string) (string, error)
}

// RedisClient is a wrapper around the Redis client.
type RedisClient struct {
	client *redis.Client
}

// NewRedisClient creates a new Redis client.
func NewRedisClient(config RedisConfig) RedisMethod {
	addr := config.Host + ":" + config.Port
	client := redis.NewClient(&redis.Options{
		Addr: addr, // Redis server address
		DB:   0,    // Default DB
	})

	// Check if the client is connected successfully
	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		log.Fatal("Error connecting to Redis:", err)
	}

	return &RedisClient{client: client}
}

// Set sets the value for the given key in Redis.
func (rc *RedisClient) Set(key string, value interface{}, expiration time.Duration) error {
	return rc.client.Set(context.Background(), key, value, expiration).Err()
}

// Get gets the value for the given key from Redis.
func (rc *RedisClient) Get(key string) (string, error) {
	return rc.client.Get(context.Background(), key).Result()
}
