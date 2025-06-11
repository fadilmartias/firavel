package controllers_v1

import (
	"goravel/app/models"
	"goravel/app/utils"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm" // TAMBAHKAN IMPORT INI
)

type AuthController struct {
	BaseController
	DB    *gorm.DB // Tambahkan ini untuk menyimpan koneksi DB
	Redis *redis.Client
}

// Ubah fungsi NewAuthController untuk menerima koneksi DB
func NewAuthController(db *gorm.DB, redis *redis.Client) *AuthController {
	return &AuthController{DB: db, Redis: redis}
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

	if result := ctrl.DB.Create(user); result.Error != nil {
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
	start := time.Now()

	// Parse input dari body
	type InputBody struct {
		Credential string `json:"credential"`
		Password   string `json:"password"`
	}
	var body InputBody
	if err := c.BodyParser(&body); err != nil {
		log.Println("⛔ BodyParser error:", err)
		log.Println("⏱️ Durasi BodyParser:", time.Since(start))
		return utils.ErrorResponse(c, utils.ErrorResponseFormat{
			Code:    fiber.StatusBadRequest,
			Message: "Invalid request body",
			Details: err.Error(),
		})
	}
	log.Println("✅ BodyParser selesai:", time.Since(start))

	// Cari user di DB berdasarkan email
	var user models.User
	dbStart := time.Now()
	if err := ctrl.DB.Where("email = ?", body.Credential).First(&user).Error; err != nil {
		log.Println("⛔ DB Query error:", err)
		log.Println("⏱️ Durasi DB Query:", time.Since(dbStart))
		log.Println("🕒 Total sampai DB:", time.Since(start))
		return utils.ErrorResponse(c, utils.ErrorResponseFormat{
			Code:    fiber.StatusUnauthorized,
			Message: "User not found",
		})
	}
	log.Println("✅ DB Query selesai:", time.Since(dbStart))

	// Bandingkan password input dengan hash di DB
	bcryptStart := time.Now()
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password)); err != nil {
		log.Println("⛔ Bcrypt compare error:", err)
		log.Println("⏱️ Durasi Bcrypt:", time.Since(bcryptStart))
		log.Println("🕒 Total sampai Bcrypt:", time.Since(start))
		return utils.ErrorResponse(c, utils.ErrorResponseFormat{
			Code:    fiber.StatusUnauthorized,
			Message: "Invalid email or password",
		})
	}
	log.Println("✅ Bcrypt compare selesai:", time.Since(bcryptStart))

	// Buat token (JWT misalnya)
	jwtStart := time.Now()
	token, err := utils.GenerateToken(&user)
	if err != nil {
		log.Println("⛔ Token gen error:", err)
		log.Println("⏱️ Durasi Token Gen:", time.Since(jwtStart))
		log.Println("🕒 Total sampai Token:", time.Since(start))
		return utils.ErrorResponse(c, utils.ErrorResponseFormat{
			Code:    fiber.StatusInternalServerError,
			Message: "Failed to generate token",
		})
	}
	log.Println("✅ Token gen selesai:", time.Since(jwtStart))

	// Logging total durasi
	log.Println("✅ Login berhasil | Total durasi login:", time.Since(start))

	return utils.SuccessResponse(c, utils.SuccessResponseFormat{
		Code:    fiber.StatusOK,
		Message: "User logged in successfully",
		Data:    fiber.Map{"token": token},
	})
}
