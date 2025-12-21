package middleware

import (
	"os"
	"stock_backend/model/response"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
)

func LoggedOutMiddleware(c *fiber.Ctx) error {
	auth := c.Get("Authorization")

	// If the token is empty, allow the request to proceed
	if auth == "" {
		return c.Next()
	}

	parts := strings.SplitN(auth, " ", 2)
	if len(parts) != 2 || parts[0] != "Bearer" {
		return c.Next()
	}
	tokenStr := parts[1]

	// Parse and validate the JWT token
	jwtSecret := os.Getenv("JWT_SECRET")
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		// Check if the signing method is HMAC
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fiber.ErrUnauthorized
		}
		return []byte(jwtSecret), nil
	})

	// If there is an error or the token is invalid, allow the request to proceed
	if err != nil || !token.Valid {
		return c.Next()
	}

	claims := token.Claims.(jwt.MapClaims)

	// Check for expired date JWT
	exp, ok := claims["exp"].(float64)
	if ok && int64(exp) > time.Now().Unix() {
		return c.Status(fiber.StatusBadRequest).JSON(response.FailedResponse{
			Message: "You are already logged in",
		})
	}

	return c.Next()
}
