package controllers_v1

import (
	"os"
	"time"

	"github.com/fadilmartias/firavel/app/mail"
	"github.com/fadilmartias/firavel/app/models"
	"github.com/fadilmartias/firavel/app/requests"
	"github.com/fadilmartias/firavel/app/utils"

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
	input, err := utils.GetValidatedBody[requests.RegisterInput](c)
	if err != nil {
		return err
	}

	bytes, err := bcrypt.GenerateFromPassword([]byte(input.Password), 14)
	if err != nil {
		return utils.ErrorResponse(c, utils.ErrorResponseFormat{
			Code:       fiber.StatusInternalServerError,
			Message:    "Could not process password",
			DevMessage: err.Error(),
			Details:    err,
		})
	}
	input.Password = string(bytes)

	if result := ctrl.DB.Create(&models.User{
		Name:     input.Name,
		Email:    input.Email,
		Password: input.Password,
	}); result.Error != nil {
		return utils.ErrorResponse(c, utils.ErrorResponseFormat{
			Code:       fiber.StatusConflict,
			Message:    "Could not create user",
			DevMessage: result.Error.Error(),
			Details:    result.Error,
		})
	}

	return utils.SuccessResponse(c, utils.SuccessResponseFormat{
		Code:    fiber.StatusCreated,
		Message: "User created successfully",
	})
}

type userResponse struct {
	ID              string     `json:"id"`
	Name            string     `json:"name"`
	Email           string     `json:"email"`
	Role            string     `json:"role"`
	EmailVerifiedAt *time.Time `json:"email_verified_at,omitempty"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

type loginResponse struct {
	User        userResponse `json:"user"`
	AccessToken string       `json:"access_token"`
}

// Login placeholder
func (ctrl *AuthController) Login(c *fiber.Ctx) error {

	// Parse input dari body
	input, err := utils.GetValidatedBody[requests.LoginInput](c)
	if err != nil {
		return err
	}

	// Cari user di DB berdasarkan email
	var user models.User
	if err := ctrl.DB.Where("email = ?", input.Credential).First(&user).Error; err != nil {
		return utils.ErrorResponse(c, utils.ErrorResponseFormat{
			Code:       fiber.StatusUnauthorized,
			Message:    "User not found",
			DevMessage: err.Error(),
			Details:    err,
		})
	}

	// Bandingkan password input dengan hash di DB
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		return utils.ErrorResponse(c, utils.ErrorResponseFormat{
			Code:       fiber.StatusUnauthorized,
			Message:    "Invalid email or password",
			DevMessage: err.Error(),
			Details:    err,
		})
	}

	// Buat token access dan refresh token
	accessToken, err := utils.GenerateToken(map[string]any{
		"user_id": user.ID,
		"name":    user.Name,
		"email":   user.Email,
		"role":    user.Role,
	}, time.Hour*1)
	if err != nil {
		return utils.ErrorResponse(c, utils.ErrorResponseFormat{
			Code:       fiber.StatusInternalServerError,
			Message:    "Failed to generate access token",
			DevMessage: err.Error(),
			Details:    err,
		})
	}

	refreshToken, err := utils.GenerateToken(map[string]any{
		"user_id": user.ID,
		"name":    user.Name,
		"email":   user.Email,
		"role":    user.Role,
	}, time.Hour*24)
	if err != nil {
		return utils.ErrorResponse(c, utils.ErrorResponseFormat{
			Code:       fiber.StatusInternalServerError,
			Message:    "Failed to generate refresh token",
			DevMessage: err.Error(),
			Details:    err,
		})
	}

	refreshCookie := fiber.Cookie{
		Name:     "refreshToken",
		Value:    refreshToken,
		Expires:  time.Now().Add(time.Hour * 24),
		HTTPOnly: true,
		SameSite: fiber.CookieSameSiteLaxMode,
		Secure:   os.Getenv("APP_ENV") == "production",
	}
	c.Cookie(&refreshCookie)

	user.RefreshToken = &refreshToken
	if err := ctrl.DB.Save(&user).Error; err != nil {
		return utils.ErrorResponse(c, utils.ErrorResponseFormat{
			Code:       fiber.StatusInternalServerError,
			Message:    "Failed to save refresh token",
			DevMessage: err.Error(),
			Details:    err,
		})
	}

	response := loginResponse{
		User: userResponse{
			ID:              user.ID,
			Name:            user.Name,
			Email:           user.Email,
			Role:            user.Role,
			EmailVerifiedAt: user.EmailVerifiedAt,
			CreatedAt:       user.CreatedAt,
			UpdatedAt:       user.UpdatedAt,
		},
		AccessToken: accessToken,
	}

	return utils.SuccessResponse(c, utils.SuccessResponseFormat{
		Code:    fiber.StatusOK,
		Message: "User logged in successfully",
		Data:    response,
	})
}

func (ctrl *AuthController) ForgotPassword(c *fiber.Ctx) error {
	input, err := utils.GetValidatedBody[requests.ForgotPasswordInput](c)
	if err != nil {
		return err
	}

	var user = models.User{}
	if err := ctrl.DB.Where("email = ?", input.Email).First(&user).Error; err != nil {
		return utils.ErrorResponse(c, utils.ErrorResponseFormat{
			Code:       fiber.StatusNotFound,
			Message:    "User not found",
			DevMessage: err.Error(),
			Details:    err,
		})
	}

	// Generate random token
	token := utils.GenerateShortID(32)
	passwordResetToken := models.PasswordResetToken{
		Email:     user.Email,
		Token:     token,
		ExpiredAt: time.Now().Add(time.Minute * 5),
	}
	passwordResetToken.HashToken(&token)
	ctrl.DB.Create(&passwordResetToken)

	// Send email
	if err := mail.SendResetPasswordEmail(user.Email, token); err != nil {
		return utils.ErrorResponse(c, utils.ErrorResponseFormat{
			Code:       fiber.StatusInternalServerError,
			Message:    "Failed to send reset password email",
			DevMessage: err.Error(),
			Details:    err,
		})
	}

	return utils.SuccessResponse(c, utils.SuccessResponseFormat{
		Code:    fiber.StatusOK,
		Message: "Password reset token sent successfully",
	})
}

func (ctrl *AuthController) ResetPassword(c *fiber.Ctx) error {
	input, err := utils.GetValidatedBody[requests.ResetPasswordInput](c)
	if err != nil {
		return err
	}

	var passwordResetToken models.PasswordResetToken
	if err := ctrl.DB.Where("token = ?", input.Token).First(&passwordResetToken).Error; err != nil {
		return utils.ErrorResponse(c, utils.ErrorResponseFormat{
			Code:       fiber.StatusNotFound,
			Message:    "Password reset token not found",
			DevMessage: err.Error(),
			Details:    err,
		})
	}

	if passwordResetToken.ExpiredAt.Before(time.Now()) {
		return utils.ErrorResponse(c, utils.ErrorResponseFormat{
			Code:       fiber.StatusNotFound,
			Message:    "Password reset token expired",
			DevMessage: "Password reset token expired",
			Details:    nil,
		})
	}

	user := models.User{}
	if err := ctrl.DB.Where("email = ?", passwordResetToken.Email).First(&user).Error; err != nil {
		return utils.ErrorResponse(c, utils.ErrorResponseFormat{
			Code:       fiber.StatusNotFound,
			Message:    "User not found",
			DevMessage: err.Error(),
			Details:    err,
		})
	}

	user.HashPassword(input.Password)
	if err := ctrl.DB.Save(&user).Error; err != nil {
		return utils.ErrorResponse(c, utils.ErrorResponseFormat{
			Code:       fiber.StatusInternalServerError,
			Message:    "Failed to save user",
			DevMessage: err.Error(),
			Details:    err,
		})
	}

	return utils.SuccessResponse(c, utils.SuccessResponseFormat{
		Code:    fiber.StatusOK,
		Message: "Password reset successfully",
	})
}
