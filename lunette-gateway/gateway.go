package main

import (
	"lunette-gateway/middleware"
	"os"

	"github.com/gin-gonic/gin"
)

func InitGateway(gateway *gin.Engine) {
	port := os.Getenv("GATEWAY_PORT")
	gateway.Run(":" + port)

	middleware.InitMiddleWares(gateway)
	router.SetupRoutingGroups(gateway)
}
