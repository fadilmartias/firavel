package middleware

import (
	"github.com/fadilmartias/firavel/app/utils"
	"github.com/gofiber/fiber/v2"
)

func Guest() fiber.Handler {
	return func(c *fiber.Ctx) error {
		if c.Locals("user") != nil || c.Get("Authorization") != "" {
			return utils.ErrorResponse(c, utils.ErrorResponseFormat{
				Code:    fiber.StatusForbidden,
				Message: "Forbidden: User already logged in",
			})
		}
		return c.Next()
	}
}
