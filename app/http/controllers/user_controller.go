package controllers

import (
	"goravel/app/models"
	"goravel/app/utils"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm" // TAMBAHKAN IMPORT INI
)

type UserController struct {
	BaseController
	DB *gorm.DB // Tambahkan ini untuk menyimpan koneksi DB
}

// Ubah fungsi NewUserController untuk menerima koneksi DB
func NewUserController(db *gorm.DB) *UserController {
	return &UserController{DB: db}
}

// Index mengambil semua user
func (ctrl *UserController) Index(c *fiber.Ctx) error {
	// Gunakan koneksi DB yang sudah ada di struct
	var users []models.User
	ctrl.DB.Find(&users)
	return utils.SuccessResponse(c, utils.SuccessResponseFormat{
		Code:    fiber.StatusOK,
		Message: "User fetched successfully",
		Data:    users,
	})
}

// Show mengambil satu user
func (ctrl *UserController) Show(c *fiber.Ctx) error {
	id := c.Params("id")
	var user models.User
	if err := ctrl.DB.First(&user, id).Error; err != nil {
		return utils.ErrorResponse(c, utils.ErrorResponseFormat{
			Code:    fiber.StatusNotFound,
			Message: "User not found",
			Details: nil,
		})
	}
	return utils.SuccessResponse(c, utils.SuccessResponseFormat{
		Code:    fiber.StatusOK,
		Message: "User fetched successfully",
		Data:    user,
	})
}

// Register membuat user baru
func (ctrl *UserController) Register(c *fiber.Ctx) error {
	user := new(models.User)

	if err := c.BodyParser(user); err != nil {
		return utils.ErrorResponse(c, utils.ErrorResponseFormat{
			Code:    fiber.StatusBadRequest,
			Message: "Invalid request body",
			Details: err.Error(),
		})
	}

	if err := user.HashPassword(user.Password); err != nil {
		return utils.ErrorResponse(c, utils.ErrorResponseFormat{
			Code:    fiber.StatusInternalServerError,
			Message: "Could not process password",
			Details: nil,
		})
	}

	if result := ctrl.DB.Create(&user); result.Error != nil {
		return utils.ErrorResponse(c, utils.ErrorResponseFormat{
			Code:    fiber.StatusConflict,
			Message: "Could not create user",
			Details: result.Error.Error(),
		})
	}

	return utils.SuccessResponse(c, utils.SuccessResponseFormat{
		Code:    fiber.StatusCreated,
		Message: "User created successfully",
		Data:    user,
	})
}

// Login placeholder
func (ctrl *UserController) Login(c *fiber.Ctx) error {
	return utils.SuccessResponse(c, utils.SuccessResponseFormat{
		Code:    fiber.StatusOK,
		Message: "User logged in successfully",
		Data:    fiber.Map{"token": "valid-token"},
	})
}
