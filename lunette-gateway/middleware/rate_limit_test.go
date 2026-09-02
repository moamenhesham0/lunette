package middleware

import (
	"lunette-gateway/util"
	"net/http"
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

var testRateLimitRequestLimit = 60

func mockRedisSetup() *redis.Client {
	miniRedis, _ := miniredis.Run()
	return redis.NewClient(&redis.Options{Addr: miniRedis.Addr()})
}

func rateLimitSetup(chain gin.HandlersChain) *gin.Engine {
	gin.SetMode(gin.TestMode)
	gateway := gin.New()
	gateway.Use(chain...)
	gateway.GET("/test", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Test successful"})
	})
	return gateway
}

func TestRateLimit_NoUserIDProvided(t *testing.T) {
	// Given
	gateway := rateLimitSetup(gin.HandlersChain{
		RateLimit(mockRedisSetup(), testRateLimitRequestLimit),
	})

	// When
	w := util.SendTestRequests(gateway, "GET", "/test")

	// Then
	if w.Code != http.StatusOK {
		expectationLog := util.UnitTestExpectation(w.Code, http.StatusOK)
		t.Errorf("Rate Limiter Should handle requests without a user ID. %s", expectationLog)
	}
}

func TestRateLimit_UserIDProvided(t *testing.T) {
	// Given
	gateway := rateLimitSetup(gin.HandlersChain{
		func(c *gin.Context) {
			c.Set("user_id", "123e4567-e89b-12d3-a456-426614174000")
		},
		RateLimit(mockRedisSetup(), testRateLimitRequestLimit),
	})

	// When
	w := util.SendTestRequests(gateway, "GET", "/test")

	// Then
	if w.Code != http.StatusOK {
		expectationLog := util.UnitTestExpectation(w.Code, http.StatusOK)
		t.Errorf("Rate Limiter Should handle requests with a user ID. %s", expectationLog)
	}
}

func TestRateLimit_RequestLimitExceeded(t *testing.T) {
	// Given
	gateway := rateLimitSetup(gin.HandlersChain{
		RateLimit(mockRedisSetup(), testRateLimitRequestLimit),
	})

	// When
	for i := 0; i < testRateLimitRequestLimit; i++ {
		w := util.SendTestRequests(gateway, "GET", "/test")

		if w.Code != http.StatusOK {
			expectationLog := util.UnitTestExpectation(w.Code, http.StatusOK)
			t.Errorf(
				"Rate Limiter Should handle %d requests within the limit but handled %d. %s",
				testRateLimitRequestLimit,
				i,
				expectationLog,
			)
		}
	}

	w := util.SendTestRequests(gateway, "GET", "/test")

	if w.Code != http.StatusTooManyRequests {
		expectationLog := util.UnitTestExpectation(w.Code, http.StatusTooManyRequests)
		t.Errorf("Rate Limiter Should reject requests that exceed the limit. %s", expectationLog)
	}
}
