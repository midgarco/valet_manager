package manager

import (
	"errors"
	"fmt"
	"time"

	"github.com/go-redis/redis"
	"github.com/midgarco/valet_manager/pkg/config"
)

// Connection ...
type Connection struct {
	Redis *redis.Client
}

// RedisConnection establishes the connection to the Reids server
func (c *Connection) RedisConnection() error {
	host := config.Get("REDIS_HOST")
	port := config.GetWithDefault("REDIS_PORT", "6379")

	if host == "" {
		return errors.New("no redis connection configured")
	}

	c.Redis = redis.NewClient(&redis.Options{
		Addr:        fmt.Sprintf("%s:%s", host, port),
		Password:    "",
		DB:          0,
		MaxRetries:  2,
		PoolTimeout: 1 * time.Second,
		IdleTimeout: 1 * time.Second,
	})

	_, err := c.Redis.Ping().Result()
	if err != nil {
		return err
	}
	return nil
}
