package middlewares

import (
	"fmt"

	"github.com/fadilmartias/firavel/app/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
)

func Auth(allowedRoles []string, allowedPermissions []string) fiber.Handler {
	return func(c *fiber.Ctx) error {

		token := c.Get("Authorization")

		if token == "" {
			return utils.ErrorResponse(c, utils.ErrorResponseFormat{
				Code:    fiber.StatusUnauthorized,
				Message: "Unauthorized: Missing token",
			})
		}

		token = token[7:] // Remove "Bearer " prefix
		var claims jwt.MapClaims
		var err error
		// Bypass token
		if token == "TestAdmin123" {
			fmt.Println("✅ Bypass token")
			// Simpan user langsung ke context
			claims = jwt.MapClaims{
				"id":    "00000A1",
				"name":  "Admin",
				"email": "admin@gmail.com",
				"phone": "08123456789",
				"role":  "admin",
			}
			return c.Next()
		} else if token == "TestUser123" {
			fmt.Println("✅ Bypass token")
			// Simpan user langsung ke context
			claims = jwt.MapClaims{
				"id":    "00000A2",
				"name":  "User",
				"email": "user@gmail.com",
				"phone": "08123456789",
				"role":  "user",
			}
		} else {
			claims, err = utils.ValidateToken(token)
		}

		if err != nil {
			return utils.ErrorResponse(c, utils.ErrorResponseFormat{
				Code:    fiber.StatusUnauthorized,
				Message: "Unauthorized: Invalid token",
			})
		}

		if claims["id"] == nil {
			return utils.ErrorResponse(c, utils.ErrorResponseFormat{
				Code:    fiber.StatusUnauthorized,
				Message: "Unauthorized: Invalid token",
			})
		}

		userID := claims["id"].(string)
		if userID == "" {
			return utils.ErrorResponse(c, utils.ErrorResponseFormat{
				Code:    fiber.StatusUnauthorized,
				Message: "Unauthorized: Invalid token",
			})
		}
		if len(allowedRoles) > 0 {
			role, ok := claims["role"].(string)
			if !ok || !utils.SliceContains(allowedRoles, role) {
				return utils.ErrorResponse(c, utils.ErrorResponseFormat{
					Code:    fiber.StatusForbidden,
					Message: "Forbidden: Access denied",
				})
			}
		}

		c.Locals("user", claims)
		return c.Next()
	}
}
