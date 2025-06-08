package routes

import (
	controllers_v1 "goravel/app/http/controllers/v1"
	"goravel/app/http/middleware"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func RegisterApiRoutes(app *fiber.App, db *gorm.DB) {

	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "Hello, World!",
		})
	})
	apiV1 := app.Group("/v1")
	userController := controllers_v1.NewUserController(db)
	authController := controllers_v1.NewAuthController(db)

	authRoutes := apiV1.Group("/auth")
	{
		authRoutes.Post("/login", authController.Login).Name("auth.login")
		authRoutes.Post("/register", authController.Register).Name("auth.register")
	}

	userRoutes := apiV1.Group("/users", middleware.Auth())
	{
		userRoutes.Get("/", userController.Index).Name("users.index")
		userRoutes.Get("/:id", userController.Show).Name("users.show")
	}
}
