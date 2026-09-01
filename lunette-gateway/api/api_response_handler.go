package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func HandleError(
	c *gin.Context,
	code int,
	err error,
	clientMsg string,
) {
	status := http.StatusText(code)

	c.JSON(code, APIErrorResponse{
		Error: APIErrorDetail{
			Code:    code,
			Status:  status,
			Message: clientMsg,
		},
	})
}
