package middleware

import (
	"lunette-gateway/api"
	"lunette-gateway/db"
	"net/http"
	"uuid"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

var rateLimitScript = redis.NewScript(`
	local current = tonumber(redis.call("GET", KEYS[1])) or 0
	if current >= tonumber(ARGV[1]) then
		return 1
	end

	redis.call("INCR", KEYS[1])
	redis.call("EXPIRE", KEYS[1], 60)
	return 0
`)

func extractUserId(c *gin.Context) string {
	id, exists := c.Get("user_id")
	if !exists {
		return c.ClientIP()
	}
	return id.(uuid.UUID).String()
}

func RateLimit(redisClient *redis.Client, requestLimit int) gin.HandlerFunc {

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
		exitCode, err := rateLimitScript.Run(c, redisClient, []string{userId}, requestLimit).Int()
		if exitCode == 1 {
			api.HandleError(c, http.StatusTooManyRequests, err, "Rate Limit Exceeded")
			c.Abort()
			return
		}

		c.Next()
	}
}
