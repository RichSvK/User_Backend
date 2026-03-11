package middleware

import (
	"net/http"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
)

func JWTMiddleware(secretKey string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		auth := c.Get("Authorization")

		// If the token is empty, allow the request to proceed
		if auth == "" {
			return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
				"message": "Authorization header is required",
			})
		}

		parts := strings.SplitN(auth, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
				"message": "Unauthorized",
			})
		}
		tokenStr := parts[1]

		// Parse and validate the JWT token
		token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
			// Check if the signing method is HMAC
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, c.Status(http.StatusUnauthorized).JSON(fiber.Map{
					"message": "Unauthorized",
				})
			}
			return []byte(secretKey), nil
		})

		if err != nil || !token.Valid {
			return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
				"message": "Unauthorized",
			})
		}

		claims := token.Claims.(jwt.MapClaims)

		// Check for expired date JWT
		exp, ok := claims["exp"].(float64)
		if !ok || int64(exp) < time.Now().Unix() {
			return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
				"message": "Token has expired",
			})
		}

		sub, ok := claims["sub"].(string)
		if !ok {
			return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
				"message": "Invalid token subject",
			})
		}
		c.Locals("userId", sub)

		role, ok := claims["role"].(string)
		if !ok {
			return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
				"message": "Invalid token role",
			})
		}
		c.Locals("role", role)

		return c.Next()
	}
}
