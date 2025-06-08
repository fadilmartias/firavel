package controllers_v1

import (
	"goravel/app/models"
	"goravel/app/utils"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm" // TAMBAHKAN IMPORT INI
)

type AuthController struct {
	BaseController
	DB *gorm.DB // Tambahkan ini untuk menyimpan koneksi DB
}

// Ubah fungsi NewAuthController untuk menerima koneksi DB
func NewAuthController(db *gorm.DB) *AuthController {
	return &AuthController{DB: db}
}

// Register membuat user baru
func (ctrl *AuthController) Register(c *fiber.Ctx) error {
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
func (ctrl *AuthController) Login(c *fiber.Ctx) error {
	// Parse input dari body
	type InputBody struct {
		Credential string `json:"credential"`
		Password   string `json:"password"`
	}
	var body InputBody
	if err := c.BodyParser(&body); err != nil {
		return utils.ErrorResponse(c, utils.ErrorResponseFormat{
			Code:    fiber.StatusBadRequest,
			Message: "Invalid request body",
			Details: err.Error(),
		})
	}

	// Cari user di DB berdasarkan email
	var user models.User
	if err := ctrl.DB.Where("email = ?", body.Credential).First(&user).Error; err != nil {
		return utils.ErrorResponse(c, utils.ErrorResponseFormat{
			Code:    fiber.StatusUnauthorized,
			Message: "User not found",
		})
	}

	// Bandingkan password input dengan hash di DB
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password)); err != nil {
		return utils.ErrorResponse(c, utils.ErrorResponseFormat{
			Code:    fiber.StatusUnauthorized,
			Message: "Invalid email or password",
		})
	}

	// Buat token (JWT misalnya)
	token, err := utils.GenerateToken(user.Id)
	if err != nil {
		return utils.ErrorResponse(c, utils.ErrorResponseFormat{
			Code:    fiber.StatusInternalServerError,
			Message: "Failed to generate token",
		})
	}

	return utils.SuccessResponse(c, utils.SuccessResponseFormat{
		Code:    fiber.StatusOK,
		Message: "User logged in successfully",
		Data:    fiber.Map{"token": token},
	})
}
