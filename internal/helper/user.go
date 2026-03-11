package helper

import "github.com/gofiber/fiber/v2"

// Get user ID in header request or locals
func GetUserID(c *fiber.Ctx) (string, bool) {
	userId := c.Get("X-User-ID")

	// Try to check from c.Locals
	if userId == "" {
		userId, _ = c.Locals("userId").(string)
	}
	return userId, userId != ""
}
