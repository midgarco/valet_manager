package manager

import (
	"errors"
	"fmt"
	"time"

	"github.com/go-redis/redis"
	"github.com/jinzhu/gorm"
	"github.com/midgarco/env"
	"github.com/midgarco/valet_manager/pkg/valet"
)

// Connection ...
type Connection struct {
	Redis *redis.Client
	DB    *gorm.DB
}

// RedisConnection establishes the connection to the Redis server
func (c *Connection) RedisConnection() error {
	host := env.Get("REDIS_HOST")
	port := env.GetWithDefault("REDIS_PORT", "6379")

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

// DBConnection establishes the connection to the database server
func (c *Connection) DBConnection() error {
	host := env.Get("DB_HOST")
	port := env.GetWithDefault("DB_PORT", "3306")
	database := env.Get("DB_NAME")
	user := env.Get("DB_USER")
	pass := env.Get("DB_PASS")

	if host == "" || database == "" {
		return errors.New("no mysql connection configured")
	}

	var err error
	c.DB, err = gorm.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True", user, pass, host, port, database))
	if err != nil {
		return err
	}

	c.DB.AutoMigrate(
		&valet.User{},
		&valet.Shift{},
		&valet.Client{},
		&valet.Address{},
		&valet.Contact{},
		&valet.PhoneNumber{},
		&valet.Employee{},
		&valet.EmployeeAvailability{},
	)

	return nil
}
