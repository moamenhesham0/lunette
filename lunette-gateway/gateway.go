package main

import (
	"log"
	"lunette-gateway/db"
	"lunette-gateway/env"
	"lunette-gateway/middleware"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

type LunetteGateway struct {
	Engine                 *gin.Engine
	EnvironmentManger      *env.LunetteEnvironmentManager
	RateLimiterRedisClient *redis.Client
}

func (gateway LunetteGateway) initEngine() {
	gateway.Engine = gin.New()

	gateway.Engine.Use(
		middleware.Auth(gateway.EnvironmentManger.JWTSecretKey),
		middleware.RateLimit(gateway.RateLimiterRedisClient, gateway.EnvironmentManger.RateLimitRequestTokenPerMin),
	)

	//TODO: uncomment when implemented
	//router.SetupRoutingGroups(gateway.Engine)
}

// Init initializes the Lunette Gateway
func (gateway LunetteGateway) Init() {
	var err error
	gateway.initEngine()
	gateway.EnvironmentManger = env.NewLunetteEnvironmentManager()
	gateway.RateLimiterRedisClient, err = db.NewRedisClient(gateway.EnvironmentManger.RedisRateLimitURL)

	if err != nil {
		log.Fatalf("Couldn't create rate limiter Redis client :" + err.Error())
	}
}
