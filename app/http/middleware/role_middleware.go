package middleware

import (
	"github.com/fadilmartias/firavel/app/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
)

func Role(allowedRoles []string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		user := c.Locals("user")
		if user == nil {
			return utils.ErrorResponse(c, utils.ErrorResponseFormat{
				Code:    fiber.StatusUnauthorized,
				Message: "Unauthorized: No user found",
			})
		}

		claims := user.(jwt.MapClaims)

		role, ok := claims["role"].(string)
		if !ok || !utils.SliceContains(allowedRoles, role) {
			return utils.ErrorResponse(c, utils.ErrorResponseFormat{
				Code:    fiber.StatusForbidden,
				Message: "Forbidden: Access denied",
			})
		}

		return c.Next()
	}
}
