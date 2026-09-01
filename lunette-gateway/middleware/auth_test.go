package middleware

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
)

var secret = os.Getenv("AUTH_JWT_SECRET")

func setupGateway() *gin.Engine {
	gateway := gin.New()
	gateway.Use(Auth())
	gateway.GET("/test", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Test successful"})
	})
	return gateway
}

func generateJWTToken(payload string) string {
	header := `{"alg": "HS256", "typ": "JWT"}`

	headerEncoded := base64.RawURLEncoding.EncodeToString([]byte(header))
	payloadEncoded := base64.RawURLEncoding.EncodeToString([]byte(payload))

	unsigendToken := headerEncoded + "." + payloadEncoded

	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(unsigendToken))

	signature := base64.RawURLEncoding.EncodeToString(mac.Sum(nil))

	return unsigendToken + "." + signature
}

func TestAuth_EmptyToken(t *testing.T) {
	// Given
	gateway := setupGateway()
	token := ""

	// When
	req, _ := http.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	gateway.ServeHTTP(w, req)

	// Then
	if w.Code != http.StatusUnauthorized {
		t.Errorf("Empty Token must be unauthorized. Expected status code %d, but got %d", http.StatusUnauthorized, w.Code)
	}
}

func TestAuth_ModifiedToken(t *testing.T) {
	// Given
	gateway := setupGateway()
	token := generateJWTToken(`{"user_id": "123e4567-e89b-12d3-a456-426614174000"}`)

	token = token[:len(token)-1] + "X"

	// When
	req, _ := http.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	gateway.ServeHTTP(w, req)

	// Then
	if w.Code != http.StatusUnauthorized {
		t.Errorf("Modified Token must be unauthorized. Expected status code %d, but got %d", http.StatusUnauthorized, w.Code)
	}
}

func TestAuth_InvalidUUIDUserID(t *testing.T) {
	// Given
	gateway := setupGateway()
	token := generateJWTToken(`{"user_id": "invalid-uuid"}`)

	// When
	req, _ := http.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	gateway.ServeHTTP(w, req)

	// Then
	if w.Code != http.StatusUnauthorized {
		t.Errorf("Invalid UUID User ID must be unauthorized. Expected status code %d, but got %d", http.StatusUnauthorized, w.Code)
	}
}

func TestAuth_ValidToken(t *testing.T) {
	// Given
	gateway := setupGateway()
	token := generateJWTToken(`{"user_id": "123e4567-e89b-12d3-a456-426614174000"}`)

	// When
	req, _ := http.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	gateway.ServeHTTP(w, req)

	// Then
	if w.Code != http.StatusOK {
		t.Errorf("Valid Token must be authorized. Expected status code %d, but got %d", http.StatusOK, w.Code)
	}
}
