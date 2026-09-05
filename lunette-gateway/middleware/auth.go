package middleware

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"lunette-gateway/api"
	"net/http"
	"strings"
	"uuid"

	"github.com/gin-gonic/gin"
)

func extractUserID(payload string) (uuid.UUID, error) {
	payloadBytes, err := base64.RawURLEncoding.DecodeString(payload)
	if err != nil {
		return uuid.Nil(), err
	}

	var target struct {
		UserID uuid.UUID `json:"user_id"`
	}
	err = json.Unmarshal(payloadBytes, &target)
	if err != nil {
		return uuid.Nil(), err
	}

	return target.UserID, nil
}

func parseJWT(token string) ([]string, error) {
	parsedToken := strings.Split(token, ".")
	if len(parsedToken) != 3 {
		return nil, errors.New("invalid token")
	}
	return parsedToken, nil
}

func verifyJWTSignature(header, payload, signature, secret string) (string, error) {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(header + "." + payload))

	resultSignature := mac.Sum(nil)
	providedSignature, err := base64.RawURLEncoding.DecodeString(signature)

	// Verify the signature
	if err != nil {
		return "Invalid signature", errors.New("invalid base64 signature")
	}

	// Verify hashed signature matching the provided one
	if !hmac.Equal(resultSignature, providedSignature) {
		return "Modified token", errors.New("signature verification failed")
	}

	return "", nil
}

// Auth returns a JWT middleware that validates the token
func Auth(JWTsecret string) gin.HandlerFunc {

	return func(c *gin.Context) {
		token := strings.TrimPrefix(c.GetHeader("Authorization"), "Bearer ")

		// Parsed Token should have [header, payload, signature] parts
		parsedToken, err := parseJWT(token)
		if err != nil {
			api.HandleError(c, http.StatusUnauthorized, err, "Invalid token")
			c.Abort()
			return
		}

		// Signature must be valid base64 and match the expected signature
		msg, err := verifyJWTSignature(parsedToken[0], parsedToken[1], parsedToken[2], JWTsecret)
		if err != nil {
			api.HandleError(c, http.StatusUnauthorized, err, msg)
			c.Abort()
			return
		}

		// User ID should be attached in the token
		userID, err := extractUserID(parsedToken[1])
		if err != nil {
			api.HandleError(c, http.StatusUnauthorized, err, "User ID not found")
			c.Abort()
			return
		}

		c.Set("user_id", userID)
		c.Next()
	}
}
