package middleware

import (
	"fmt"
	"reflect"

	"github.com/fadilmartias/firavel/app/models"
	"github.com/fadilmartias/firavel/app/utils"
	"github.com/gofiber/fiber/v2"
)

func UserOwnership(inputColumn string, modelColumn string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Ambil user dari context
		user, ok := c.Locals("user").(models.User)
		if !ok {
			return utils.ErrorResponse(c, utils.ErrorResponseFormat{
				Code:    fiber.StatusUnauthorized,
				Message: "User not found in context",
			})
		}

		// Parse body request menjadi map
		input := make(map[string]any)
		if err := c.BodyParser(&input); err != nil {
			return utils.ErrorResponse(c, utils.ErrorResponseFormat{
				Code:    fiber.StatusBadRequest,
				Message: "Invalid request body",
			})
		}

		// Ambil nilai dari input dan konversi ke string
		inputVal, ok := input[inputColumn]
		if !ok {
			return utils.ErrorResponse(c, utils.ErrorResponseFormat{
				Code:    fiber.StatusBadRequest,
				Message: fmt.Sprintf("Field '%s' not found in request body", inputColumn),
			})
		}

		// Ambil nilai dari field user yang ditentukan (via reflection)
		val := reflect.ValueOf(user)
		if val.Kind() == reflect.Ptr {
			val = val.Elem()
		}

		field := val.FieldByName(modelColumn)
		if !field.IsValid() {
			return utils.ErrorResponse(c, utils.ErrorResponseFormat{
				Code:    fiber.StatusInternalServerError,
				Message: fmt.Sprintf("User struct does not have field '%s'", modelColumn),
			})
		}

		// Bandingkan value sebagai string (untuk fleksibilitas)
		inputStr := fmt.Sprint(inputVal)
		userStr := fmt.Sprint(field.Interface())

		if inputStr != userStr {
			return utils.ErrorResponse(c, utils.ErrorResponseFormat{
				Code:    fiber.StatusForbidden,
				Message: "Forbidden: User not authorized for this resource",
			})
		}

		return c.Next()
	}
}
