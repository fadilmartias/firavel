package middleware

import "github.com/gofiber/fiber/v2"

func RobotTag() fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Set("X-Robots-Tag", "noindex, nofollow")
		return c.Next()
	}
}
