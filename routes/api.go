package routes

import (
	controllers_v0 "github.com/fadilmartias/firavel/app/http/controllers/v0"
	controllers_v1 "github.com/fadilmartias/firavel/app/http/controllers/v1"
	"github.com/fadilmartias/firavel/app/http/middleware"
	"github.com/fadilmartias/firavel/app/jobs"
	"github.com/fadilmartias/firavel/app/logger"
	"github.com/fadilmartias/firavel/app/requests"
	"github.com/fadilmartias/firavel/app/utils"
	"github.com/fadilmartias/firavel/config"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func RegisterApiRoutes(app *fiber.App, db *gorm.DB, redis *config.RedisClient) {

	app.Use(middleware.RobotTag())
	app.Use(middleware.GetUser())
	app.Get("/", func(c *fiber.Ctx) error {
		return utils.SuccessResponse(c, utils.SuccessResponseFormat{
			Message: "Hello, World!",
		})
	})
	app.Get("/send-wa/:order_id", func(c *fiber.Ctx) error {
		orderID := c.Params("order_id")
		task, err := jobs.NewProductOrderSendWAInvoiceTask(orderID)
		if err != nil {
			logger.Error("error create task:", err)
		}
		info, err := jobs.AsynqClient.Enqueue(task)
		if err != nil {
			logger.Error("error enqueue task:", err)
		}
		logger.Infof("Task enqueued: %+v", info)
		return utils.SuccessResponse(c, utils.SuccessResponseFormat{
			Message: "Berhasil mengirim pesan",
			Data:    nil,
		})
	}).Name("send-wa")

	apiV0 := app.Group("/v0")
	genericController := controllers_v0.NewGenericController(db, redis)
	apiV0.Get("/:model", genericController.Index).Name("generic.index")
	apiV0.Get("/:model/:id", genericController.Show).Name("generic.show")
	apiV0.Post("/:model", genericController.Store, middleware.Auth([]string{"admin"}, []string{})).Name("generic.store")
	apiV0.Put("/:model/:id", genericController.Update, middleware.Auth([]string{"admin"}, []string{})).Name("generic.update")
	apiV0.Delete("/:model/:id", genericController.Destroy, middleware.Auth([]string{"admin"}, []string{})).Name("generic.destroy")

	apiV1 := app.Group("/v1")
	userController := controllers_v1.NewUserController(db, redis)
	authController := controllers_v1.NewAuthController(db, redis)

	authRoutes := apiV1.Group("/auth")
	{
		authRoutes.Post("/login", middleware.Guest(), middleware.Validator[requests.LoginInput](), authController.Login).Name("auth.login")
		authRoutes.Get("/google/redirect", middleware.Guest(), authController.GoogleRedirect).Name("auth.google.redirect")
		authRoutes.Get("/google/callback", middleware.Guest(), authController.GoogleCallback).Name("auth.google.callback")
		authRoutes.Post("/register", middleware.Guest(), middleware.Validator[requests.RegisterInput](), authController.Register).Name("auth.register")
		authRoutes.Post("/forgot-password", middleware.Guest(), middleware.Validator[requests.ForgotPasswordInput](), authController.ForgotPassword).Name("auth.forgot-password")
		authRoutes.Post("/reset-password", middleware.Guest(), middleware.Validator[requests.ResetPasswordInput](), authController.ResetPassword).Name("auth.reset-password")
		authRoutes.Post("/send-email-verification", middleware.Auth([]string{}, []string{}), authController.SendEmailVerification).Name("auth.send-email-verification")
		authRoutes.Post("/verify-email", middleware.Auth([]string{}, []string{}), middleware.Validator[requests.VerifyEmailInput](), authController.VerifyEmail).Name("auth.verify-email")
	}

	userRoutes := apiV1.Group("/users", middleware.Auth([]string{"admin"}, []string{}))
	{
		userRoutes.Get("/", userController.Index).Name("users.index")
		userRoutes.Get("/:id", userController.Show).Name("users.show")
	}
}
