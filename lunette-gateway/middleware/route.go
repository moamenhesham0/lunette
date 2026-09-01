package middleware

import "github.com/gin-gonic/gin"

// Routes the request to the specific service
func Route() gin.HandlerFunc {
	return func(context *gin.Context) {
		print("Routing request...")
	}
}
