package middleware

import (
	"github.com/fadilmartias/firavel/app/models"
	"github.com/fadilmartias/firavel/app/utils"
	"github.com/gofiber/fiber/v2"
)

func UserOwnership(column string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		user := c.Locals("user").(models.User)
		input := make(map[string]any)
		if err := c.BodyParser(&input); err != nil {
			return utils.ErrorResponse(c, utils.ErrorResponseFormat{
				Code:    fiber.StatusBadRequest,
				Message: "Invalid request body",
			})
		}

		if input[column].(string) != user.ID {
			return utils.ErrorResponse(c, utils.ErrorResponseFormat{
				Code:    fiber.StatusForbidden,
				Message: "Forbidden: User not authorized",
			})
		}
		return c.Next()
	}
}
