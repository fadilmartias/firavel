package controllers_v1

import (
	"goravel/app/models"
	"goravel/app/utils"

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
