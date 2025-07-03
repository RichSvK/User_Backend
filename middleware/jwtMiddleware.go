package middleware

import (
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
)

func JWTMiddleware(c *fiber.Ctx) error {
	tokenStr := c.Cookies("token")

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

	// Check for expired date JWT
	exp, ok := claims["exp"].(float64)
	if !ok || int64(exp) < time.Now().Unix() {
		return fiber.ErrUnauthorized
	}

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
