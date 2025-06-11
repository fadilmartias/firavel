package middleware

import (
	"fmt"
	"goravel/app/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
)

func Auth() fiber.Handler {
	return func(c *fiber.Ctx) error {
		token := c.Get("Authorization")

		if token == "" {
			return utils.ErrorResponse(c, utils.ErrorResponseFormat{
				Code:    fiber.StatusUnauthorized,
				Message: "Unauthorized: Missing token",
			})
		}

		token = token[7:] // Remove "Bearer " prefix

		// Bypass token
		if token == "TestAdmin123" {
			fmt.Println("✅ Bypass token")
			// Simpan user langsung ke context
			c.Locals("user", jwt.MapClaims{
				"user_id": "00000A1",
				"name":    "Admin",
				"email":   "admin@gmail.com",
				"role":    "admin",
			})
			return c.Next()
		} else if token == "TestUser123" {
			fmt.Println("✅ Bypass token")
			// Simpan user langsung ke context
			c.Locals("user", jwt.MapClaims{
				"user_id": "00000A2",
				"name":    "User",
				"email":   "user@gmail.com",
				"role":    "user",
			})
			return c.Next()
		}

		claims, err := utils.ValidateToken(token)
		if err != nil {
			return utils.ErrorResponse(c, utils.ErrorResponseFormat{
				Code:    fiber.StatusUnauthorized,
				Message: "Unauthorized: Invalid token",
			})
		}

		if claims["user_id"] == nil {
			return utils.ErrorResponse(c, utils.ErrorResponseFormat{
				Code:    fiber.StatusUnauthorized,
				Message: "Unauthorized: Invalid token",
			})
		}

		userID := claims["user_id"].(string)
		if userID == "" {
			return utils.ErrorResponse(c, utils.ErrorResponseFormat{
				Code:    fiber.StatusUnauthorized,
				Message: "Unauthorized: Invalid token",
			})
		}
		c.Locals("user", claims)
		return c.Next()
	}
}
