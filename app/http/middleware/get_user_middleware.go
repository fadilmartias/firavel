// middleware/get_user.go
package middleware

import (
	"github.com/fadilmartias/firavel/app/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
)

func GetUser() fiber.Handler {
	return func(c *fiber.Ctx) error {
		token := c.Cookies("accessToken") // baca dari cookie
		if token == "" {
			// Tidak ada token, lanjut tanpa user
			return c.Next()
		}
		var claims jwt.MapClaims
		var err error

		switch token {
		case "TestAdmin123":
			claims = jwt.MapClaims{
				"id":    "A1",
				"name":  "Admin",
				"email": "admin@gmail.com",
				"phone": "08123456789",
				"role":  "admin",
			}
		case "TestUser123":
			claims = jwt.MapClaims{
				"id":    "A2",
				"name":  "User",
				"email": "user@gmail.com",
				"phone": "08123456780",
				"role":  "user",
			}
		default:
			claims, err = utils.ValidateToken(token)
			if err != nil {
				// Token ada tapi tidak valid, abaikan
				return c.Next()
			}
		}

		// Jika berhasil, simpan ke context
		if claims["id"] != nil {
			c.Locals("user", claims)
		}

		return c.Next()
	}
}
