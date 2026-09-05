package util

import (
	"fmt"
	"net/http/httptest"

	"github.com/gin-gonic/gin"
)

func UnitTestExpectation(got any, expected any) string {
	return fmt.Sprintf("Expected %v , got %v", expected, got)
}

func SendTestRequests(engine *gin.Engine, method string, path string, headers ...map[string]string) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	req := httptest.NewRequest(method, path, nil)
	for _, header := range headers {
		for key, value := range header {
			req.Header.Set(key, value)
		}
	}

	engine.ServeHTTP(w, req)

	return w
}
