// middleware/auth.go
package middleware

import (
	"github.com/fadilmartias/firavel/app/utils"
	"github.com/gofiber/fiber/v2"
)

func Auth(allowedRoles []string, allowedPermissions []string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userClaims := c.Locals("user")
		if userClaims == nil {
			return utils.ErrorResponse(c, utils.ErrorResponseFormat{
				Code:    fiber.StatusUnauthorized,
				Message: "Unauthorized: User not found",
			})
		}

		claims, ok := userClaims.(map[string]interface{})
		if !ok {
			return utils.ErrorResponse(c, utils.ErrorResponseFormat{
				Code:    fiber.StatusUnauthorized,
				Message: "Unauthorized: Invalid user data",
			})
		}

		// Validate role
		if len(allowedRoles) > 0 {
			role, ok := claims["role"].(string)
			if !ok || !utils.SliceContains(allowedRoles, role) {
				return utils.ErrorResponse(c, utils.ErrorResponseFormat{
					Code:    fiber.StatusForbidden,
					Message: "Forbidden: Access denied",
				})
			}
		}

		// Tambahan: validate permission kalau dibutuhkan
		// if len(allowedPermissions) > 0 {
		//    ...
		// }

		return c.Next()
	}
}
