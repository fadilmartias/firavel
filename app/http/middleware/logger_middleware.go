package middleware

import (
	config "github.com/fadilmartias/firavel/app/logger"
	"github.com/gofiber/fiber/v2"
)

func Logger() fiber.Handler {
	return func(c *fiber.Ctx) error {
		err := c.Next()

		if err != nil {
			config.Error("Fiber error: %v | Path: %s | IP: %s", err, c.Path(), c.IP())
			// Fiber default error handler
			return fiber.DefaultErrorHandler(c, err)
		}
		return nil
	}
}
