package db

import (
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func NewRedisClient(redisURL string) (*redis.Client, error) {
	opt, err := redis.ParseURL(redisURL)
	if err != nil {
		return nil, err
	}
	return redis.NewClient(opt), nil
}

func HealthCheck(redisClient *redis.Client, c *gin.Context) error {
	_, err := redisClient.Ping(c).Result()
	return err
}
