package middleware

import (
	"github.com/fadilmartias/firavel/app/utils"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

var validate = validator.New()

func ValidateBody[T any]() fiber.Handler {
	return func(c *fiber.Ctx) error {
		var body T

		if err := c.BodyParser(&body); err != nil {
			return utils.ErrorResponse(c, utils.ErrorResponseFormat{
				Code:       fiber.StatusBadRequest,
				Message:    "Invalid request body",
				DevMessage: err.Error(),
				Details:    err,
			})
		}

		if err := validate.Struct(body); err != nil {
			errs := make(map[string]string)
			for _, e := range err.(validator.ValidationErrors) {
				errs[e.Field()] = e.Tag()
			}
			return utils.ErrorResponse(c, utils.ErrorResponseFormat{
				Code:       fiber.StatusUnprocessableEntity,
				Message:    "Invalid request body",
				DevMessage: err.Error(),
				Details: map[string]interface{}{
					"errors": errs,
				},
			})
		}

		c.Locals("validatedBody", body)

		return c.Next()
	}
}
