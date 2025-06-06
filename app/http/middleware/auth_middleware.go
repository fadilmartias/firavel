package middleware

import "github.com/gofiber/fiber/v2"

// Protected adalah contoh middleware proteksi rute.
// Di aplikasi nyata, ini akan memverifikasi token JWT.
func Auth() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Contoh sederhana: cek header Authorization
		token := c.Get("Authorization")

		// Di sini Anda akan memvalidasi token (misalnya JWT)
		// Untuk contoh ini, kita hanya cek jika token ada dan valid (dummy check)
		if token != "Bearer valid-token" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"status":  "error",
				"message": "Unauthorized: Missing or invalid token",
			})
		}

		// Jika valid, lanjutkan ke handler berikutnya
		return c.Next()
	}
}