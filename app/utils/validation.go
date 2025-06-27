package utils

import (
	"github.com/gofiber/fiber/v2"
)

func GetValidatedBody[T any](c *fiber.Ctx) (T, error) {
	v := c.Locals("validatedBody")
	input, ok := v.(T)
	if !ok {
		return *new(T), ErrorResponse(c, ErrorResponseFormat{
			Code:    fiber.StatusInternalServerError,
			Message: "Failed to parse validated body",
		})
	}
	return input, nil
}
