package logger

import (
	"github.com/gofiber/fiber/v2"
)

func LoggerMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		err := c.Next()

		if err != nil {
			Errorf("Fiber error: %v | Path: %s | IP: %s", err, c.Path(), c.IP())
			// Fiber default error handler
			return fiber.DefaultErrorHandler(c, err)
		}
		return nil
	}
}
