package middlewares

import (
	"os"

	"github.com/gofiber/fiber/v2"
)

func CheckAppKey(c *fiber.Ctx) error {
	requestedAppKey := c.Get("app-key")
	appKey := os.Getenv("APP_KEY")

	if requestedAppKey != appKey {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Unauthorized",
		})
	}

	return c.Next()
}