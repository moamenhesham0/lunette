package middleware

import (
	"log"
	"lunette-gateway/api"
	"lunette-gateway/db"
	"net/http"
	"os"
	"uuid"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

var rateLimitScript = redis.NewScript(`
	local current = tonumber(redis.call("GET", KEYS[1])) or 0
	if current >= tonumber(ARGV[1]) then
		return 1
	end

	redis.call("INCR", key)
	redis.call("EXPIRE", key, 60)
	return 0
`)

func extractUserId(c *gin.Context) string {
	id, exists := c.Get("user_id")
	if !exists {
		return c.ClientIP()
	}
	return id.(uuid.UUID).String()
}

func RateLimit() gin.HandlerFunc {
	redisURL := os.Getenv("GATEWAY_RATELIMIT_REDIS_URL")
	redisClient, err := db.NewRedisClient(redisURL)

	if err != nil {
		log.Fatalf("Rate limit Redis client failed to start: %v", err)
	}

	rateLimitRequestTokensPerMinute := 60

	return func(c *gin.Context) {

		// Redis client should be healthy
		err := db.HealthCheck(redisClient, c)
		if err != nil {
			api.HandleError(c, http.StatusServiceUnavailable, err, "Internal Server Error")
			c.Abort()
			return
		}

		userId := extractUserId(c)

		// User should be within request limits
		err = rateLimitScript.Run(c, redisClient, []string{userId}, rateLimitRequestTokensPerMinute).Err()
		if err != nil {
			api.HandleError(c, http.StatusTooManyRequests, err, "Rate Limit Exceeded")
			c.Abort()
			return
		}

		c.Next()
	}
}
