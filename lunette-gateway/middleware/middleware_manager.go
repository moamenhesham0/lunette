package middleware

import (
	"github.com/gin-gonic/gin"
)

// InitMiddleWares initializes the middleware for the gateway
func InitMiddleWares(engine *gin.Engine) {
	engine.Use(
	// MiddleWare Layers
	//RateLimit(),
	//Auth(),
	//Route(),
	)
}
