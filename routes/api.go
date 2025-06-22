package routes

import (
	controllers_v0 "github.com/fadilmartias/firavel/app/http/controllers/v0"
	controllers_v1 "github.com/fadilmartias/firavel/app/http/controllers/v1"
	"github.com/fadilmartias/firavel/app/http/middleware"

	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func RegisterApiRoutes(app *fiber.App, db *gorm.DB, redis *redis.Client) {

	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "Hello, World!",
		})
	})
	apiV0 := app.Group("/v0")
	genericController := controllers_v0.NewGenericController(db, redis)
	apiV0.Get("/:model", genericController.Index).Name("generic.index")
	apiV0.Get("/:model/:id", genericController.Show).Name("generic.show")
	apiV0.Post("/:model", genericController.Store).Name("generic.store")
	apiV0.Put("/:model/:id", genericController.Update).Name("generic.update")
	apiV0.Delete("/:model/:id", genericController.Destroy).Name("generic.destroy")

	apiV1 := app.Group("/v1")
	userController := controllers_v1.NewUserController(db, redis)
	authController := controllers_v1.NewAuthController(db, redis)

	authRoutes := apiV1.Group("/auth")
	{
		authRoutes.Post("/login", authController.Login).Name("auth.login")
		authRoutes.Post("/register", authController.Register).Name("auth.register")
	}

	userRoutes := apiV1.Group("/users", middleware.Auth(), middleware.Role([]string{"admin"}))
	{
		userRoutes.Get("/", userController.Index).Name("users.index")
		userRoutes.Get("/:id", userController.Show).Name("users.show")
	}
}
