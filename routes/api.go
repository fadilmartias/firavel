package routes

import (
	"goravel/app/http/controllers"
	"goravel/app/http/middleware"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func RegisterApiRoutes(app *fiber.App, db *gorm.DB) {
	userController := controllers.NewUserController(db)
	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "Hello, World!",
		})
	})
	api := app.Group("/v1")

	api.Post("/login", userController.Login).Name("auth.login")
	api.Post("/register", userController.Register).Name("auth.register")

	authRoutes := api.Group("/users", middleware.Auth())
	{
		authRoutes.Get("/", userController.Index).Name("users.index")
		authRoutes.Get("/:id", userController.Show).Name("users.show")
	}
}
