package middleware

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"lunette-gateway/util"
	"net/http"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

var testAuthSecret = "SECRET-KEY0-XYZ"

func authSetup() *gin.Engine {
	gin.SetMode(gin.TestMode)
	gateway := gin.New()
	gateway.Use(Auth(testAuthSecret))
	gateway.GET("/test", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Test successful"})
	})
	return gateway
}

func modifyTokenPayload(token string, modifiedPayload string) string {
	modifiedPayloadEncoded := base64.RawURLEncoding.EncodeToString([]byte(modifiedPayload))
	parsedToken := strings.Split(token, ".")

	return parsedToken[0] + "." + modifiedPayloadEncoded + "." + parsedToken[2]
}
func generateJWTToken(payload string) string {
	header := `{"alg": "HS256", "typ": "JWT"}`

	headerEncoded := base64.RawURLEncoding.EncodeToString([]byte(header))
	payloadEncoded := base64.RawURLEncoding.EncodeToString([]byte(payload))

	unsigendToken := headerEncoded + "." + payloadEncoded

	mac := hmac.New(sha256.New, []byte(testAuthSecret))
	mac.Write([]byte(unsigendToken))

	signature := base64.RawURLEncoding.EncodeToString(mac.Sum(nil))

	return unsigendToken + "." + signature
}

func TestAuth_EmptyToken(t *testing.T) {
	// Given
	gateway := authSetup()
	token := ""

	// When
	w := util.SendTestRequests(
		gateway,
		"GET",
		"/test",
		map[string]string{"Authorization": "Bearer " + token},
	)

	// Then
	if w.Code != http.StatusUnauthorized {
		expectationLog := util.UnitTestExpectation(w.Code, http.StatusUnauthorized)
		t.Errorf("Empty Token must be unauthorized. %s", expectationLog)
	}
}

func TestAuth_ModifiedToken(t *testing.T) {
	// Given
	gateway := authSetup()
	token := generateJWTToken(`{"user_id": "123e4567-e89b-12d3-a456-426614174000"}`)

	modifiedPayload := `{"user_id": "123e4567-e89b-12d3-a456-000000000000"`
	modifiedToken := modifyTokenPayload(token, modifiedPayload)

	// When
	w := util.SendTestRequests(
		gateway,
		"GET",
		"/test",
		map[string]string{"Authorization": "Bearer " + modifiedToken},
	)

	// Then
	if w.Code != http.StatusUnauthorized {
		expectationLog := util.UnitTestExpectation(w.Code, http.StatusUnauthorized)
		t.Errorf("Modified Token must be unauthorized. %s", expectationLog)
	}
}

func TestAuth_InvalidUUIDUserID(t *testing.T) {
	// Given
	gateway := authSetup()
	token := generateJWTToken(`{"user_id": "invalid-uuid"}`)

	// When
	w := util.SendTestRequests(
		gateway,
		"GET",
		"/test",
		map[string]string{"Authorization": "Bearer " + token},
	)

	// Then
	if w.Code != http.StatusUnauthorized {
		expectationLog := util.UnitTestExpectation(w.Code, http.StatusUnauthorized)
		t.Errorf("Invalid UUID User ID must be unauthorized. %s", expectationLog)
	}
}

func TestAuth_ValidToken(t *testing.T) {
	// Given
	gateway := authSetup()
	token := generateJWTToken(`{"user_id": "123e4567-e89b-12d3-a456-426614174000"}`)

	// When
	w := util.SendTestRequests(
		gateway,
		"GET",
		"/test",
		map[string]string{"Authorization": "Bearer " + token},
	)

	// Then
	if w.Code != http.StatusOK {
		expectationLog := util.UnitTestExpectation(w.Code, http.StatusOK)
		t.Errorf("Valid Token must be authorized. %s", expectationLog)
	}
}
