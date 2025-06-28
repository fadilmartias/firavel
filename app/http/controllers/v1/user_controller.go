package controllers_v1

import (
	"github.com/fadilmartias/firavel/app/models"
	"github.com/fadilmartias/firavel/app/requests"
	"github.com/fadilmartias/firavel/app/utils"
	"golang.org/x/crypto/bcrypt"

	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm" // TAMBAHKAN IMPORT INI
)

type UserController struct {
	BaseController
	DB    *gorm.DB // Tambahkan ini untuk menyimpan koneksi DB
	Redis *redis.Client
}

// Ubah fungsi NewUserController untuk menerima koneksi DB
func NewUserController(db *gorm.DB, redis *redis.Client) *UserController {
	return &UserController{DB: db, Redis: redis}
}

// Index mengambil semua user
func (ctrl *UserController) Index(c *fiber.Ctx) error {
	// Gunakan koneksi DB yang sudah ada di struct
	type UserResponse struct {
		ID    string `json:"id"`
		Name  string `json:"name"`
		Email string `json:"email"`
	}

	var users []UserResponse
	ctrl.DB.Model(&models.User{}).Select("id", "name", "email").Scan(&users)
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

func (ctrl *UserController) UpdateProfile(c *fiber.Ctx) error {
	input, err := utils.GetValidatedBody[requests.UpdateProfileInput](c)
	if err != nil {
		return err
	}
	id := c.Locals("user").(models.User).ID
	var user models.User
	if err := ctrl.DB.First(&user, id).Error; err != nil {
		return utils.ErrorResponse(c, utils.ErrorResponseFormat{
			Code:    fiber.StatusNotFound,
			Message: "User not found",
			Details: nil,
		})
	}
	if err := ctrl.DB.Model(&user).Updates(input).Error; err != nil {
		return utils.ErrorResponse(c, utils.ErrorResponseFormat{
			Code:    fiber.StatusInternalServerError,
			Message: "Failed to update user",
			Details: nil,
		})
	}
	return utils.SuccessResponse(c, utils.SuccessResponseFormat{
		Code:    fiber.StatusOK,
		Message: "Profile updated successfully",
		Data:    user,
	})
}

func (ctrl *UserController) UpdatePassword(c *fiber.Ctx) error {
	input, err := utils.GetValidatedBody[requests.UpdatePasswordInput](c)
	if err != nil {
		return err
	}
	id := c.Locals("user").(models.User).ID
	var user models.User
	if err := ctrl.DB.First(&user, id).Error; err != nil {
		return utils.ErrorResponse(c, utils.ErrorResponseFormat{
			Code:    fiber.StatusNotFound,
			Message: "User not found",
			Details: nil,
		})
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.CurrentPassword)); err != nil {
		return utils.ErrorResponse(c, utils.ErrorResponseFormat{
			Code:       fiber.StatusUnauthorized,
			Message:    "Invalid current password",
			DevMessage: err.Error(),
			Details:    err,
		})
	}
	bytes, err := bcrypt.GenerateFromPassword([]byte(input.NewPassword), 14)
	if err != nil {
		return utils.ErrorResponse(c, utils.ErrorResponseFormat{
			Code:       fiber.StatusInternalServerError,
			Message:    "Could not process password",
			DevMessage: err.Error(),
			Details:    err,
		})
	}
	input.NewPassword = string(bytes)
	if err := ctrl.DB.Model(&user).Updates(map[string]any{
		"password": input.NewPassword,
	}).Error; err != nil {
		return utils.ErrorResponse(c, utils.ErrorResponseFormat{
			Code:    fiber.StatusInternalServerError,
			Message: "Failed to update user",
			Details: nil,
		})
	}
	return utils.SuccessResponse(c, utils.SuccessResponseFormat{
		Code:    fiber.StatusOK,
		Message: "Password updated successfully",
		Data:    user,
	})
}
