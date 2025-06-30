package middleware

import (
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
)

func JWTMiddleware(c *fiber.Ctx) error {
	// Get value of "Authorization" header
	auth := c.Get("Authorization")
	if auth == "" {
		return fiber.ErrUnauthorized
	}

	parts := strings.SplitN(auth, " ", 2)
	if len(parts) != 2 || parts[0] != "Bearer" {
		return fiber.ErrUnauthorized
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

	if err != nil || !token.Valid {
		return fiber.ErrUnauthorized
	}

	claims := token.Claims.(jwt.MapClaims)
	sub, ok := claims["sub"].(string)
	if !ok {
		return fiber.ErrUnauthorized
	}
	c.Locals("userId", sub)

	role, ok := claims["role"].(string)
	if !ok {
		return fiber.ErrUnauthorized
	}
	c.Locals("role", role)
	
	return c.Next()
}
