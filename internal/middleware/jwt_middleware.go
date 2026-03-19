package middleware

import (
	"stock_backend/internal/handler"
	"stock_backend/internal/model/domainerr"
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
			return handler.ResponseErrorJSON(c, fiber.StatusUnauthorized, domainerr.ErrAuthorizationHeaderRequired.Error())
		}

		parts := strings.SplitN(auth, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			return handler.ResponseErrorJSON(c, fiber.StatusUnauthorized, domainerr.ErrUnauthorized.Error())
		}
		tokenStr := parts[1]

		// Parse and validate the JWT token
		token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
			// Check if the signing method is HMAC
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, handler.ResponseErrorJSON(c, fiber.StatusUnauthorized, domainerr.ErrUnauthorized.Error())
			}
			return []byte(secretKey), nil
		})

		if err != nil || !token.Valid {
			return handler.ResponseErrorJSON(c, fiber.StatusUnauthorized, domainerr.ErrUnauthorized.Error())
		}

		claims := token.Claims.(jwt.MapClaims)

		// Check for expired date JWT
		exp, ok := claims["exp"].(float64)
		if !ok || int64(exp) < time.Now().Unix() {
			return handler.ResponseErrorJSON(c, fiber.StatusUnauthorized, domainerr.ErrTokenExpired.Error())
		}

		sub, ok := claims["sub"].(string)
		if !ok {
			return handler.ResponseErrorJSON(c, fiber.StatusUnauthorized, domainerr.ErrInvalidTokenClaims.Error())
		}
		c.Locals("userId", sub)

		role, ok := claims["role"].(string)
		if !ok {
			return handler.ResponseErrorJSON(c, fiber.StatusUnauthorized, domainerr.ErrInvalidTokenClaims.Error())
		}
		c.Locals("role", role)

		return c.Next()
	}
}
