package controllers_v1

import (
	"fmt"
	"os"
	"time"

	"github.com/fadilmartias/firavel/app/mail"
	"github.com/fadilmartias/firavel/app/models"
	"github.com/fadilmartias/firavel/app/requests"
	"github.com/fadilmartias/firavel/app/utils"
	"github.com/golang-jwt/jwt"

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
		Phone:    input.Phone,
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
	Phone           string     `json:"phone"`
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
		"id":    user.ID,
		"name":  user.Name,
		"phone": user.Phone,
		"email": user.Email,
		"role":  user.Role,
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
		"id":    user.ID,
		"name":  user.Name,
		"phone": user.Phone,
		"email": user.Email,
		"role":  user.Role,
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
			Phone:           user.Phone,
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
	var passwordResetToken models.PasswordResetToken
	if err := ctrl.DB.Where("email = ?", input.Email).First(&user).Error; err != nil {
		return utils.ErrorResponse(c, utils.ErrorResponseFormat{
			Code:       fiber.StatusNotFound,
			Message:    "User not found",
			DevMessage: err.Error(),
			Details:    err,
		})
	}

	if err := ctrl.DB.Where("email = ?", input.Email).Where("expired_at > ?", time.Now()).First(&passwordResetToken).Error; err == nil {
		return utils.SuccessResponse(c, utils.SuccessResponseFormat{
			Code:    fiber.StatusOK,
			Message: "Password reset token already sent",
		})
	}

	// Generate random token
	token := utils.GenerateRandomToken(32)
	tokenHash := utils.HashToken(token)

	passwordResetToken.Email = user.Email
	passwordResetToken.Token = tokenHash
	passwordResetToken.ExpiredAt = time.Now().Add(time.Minute * 5)
	ctrl.DB.Save(&passwordResetToken)

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
	tokenHash := utils.HashToken(input.Token)
	var passwordResetToken models.PasswordResetToken
	if err := ctrl.DB.Where("token = ?", tokenHash).First(&passwordResetToken).Error; err != nil {
		return utils.ErrorResponse(c, utils.ErrorResponseFormat{
			Code:       fiber.StatusNotFound,
			Message:    "Password reset token not found",
			DevMessage: err.Error(),
			Details:    err,
		})
	}

	if passwordResetToken.UsedAt != nil {
		return utils.ErrorResponse(c, utils.ErrorResponseFormat{
			Code:    fiber.StatusConflict,
			Message: "Password reset token already used",
			Details: nil,
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

	now := time.Now()

	passwordResetToken.UsedAt = &now
	ctrl.DB.Save(&passwordResetToken)

	return utils.SuccessResponse(c, utils.SuccessResponseFormat{
		Code:    fiber.StatusOK,
		Message: "Password reset successfully",
	})
}

func (ctrl *AuthController) SendEmailVerification(c *fiber.Ctx) error {
	id := c.Locals("user").(jwt.MapClaims)["id"].(string)
	var user models.User
	if err := ctrl.DB.Where("id = ?", id).First(&user).Error; err != nil {
		return utils.ErrorResponse(c, utils.ErrorResponseFormat{
			Code:    fiber.StatusNotFound,
			Message: "User not found",
			Details: nil,
		})
	}
	if user.EmailVerifiedAt != nil {
		return utils.SuccessResponse(c, utils.SuccessResponseFormat{
			Code:    fiber.StatusOK,
			Message: "Email already verified",
		})
	}
	_, err := ctrl.Redis.Get(c.UserContext(), fmt.Sprintf("email_verification_token:%s", user.Email)).Result()
	if err == nil {
		return utils.ErrorResponse(c, utils.ErrorResponseFormat{
			Code:    fiber.StatusConflict,
			Message: "Email verification already sent",
		})
	}

	jwtToken, _ := utils.GenerateToken(map[string]any{
		"email": user.Email,
	}, time.Minute*5)

	if err := ctrl.Redis.SetEX(c.UserContext(), fmt.Sprintf("email_verification_token:%s", user.Email), jwtToken, time.Minute*5).Err(); err != nil {
		return utils.ErrorResponse(c, utils.ErrorResponseFormat{
			Code:       fiber.StatusInternalServerError,
			Message:    "Failed to save email verification token",
			DevMessage: err.Error(),
			Details:    err,
		})
	}

	// Send email
	if err := mail.SendEmailVerificationEmail(user.Email, jwtToken); err != nil {
		return utils.ErrorResponse(c, utils.ErrorResponseFormat{
			Code:       fiber.StatusInternalServerError,
			Message:    "Failed to send email verification",
			DevMessage: err.Error(),
			Details:    err,
		})
	}

	return utils.SuccessResponse(c, utils.SuccessResponseFormat{
		Code:    fiber.StatusOK,
		Message: "Email verification sent successfully",
	})
}

func (ctrl *AuthController) VerifyEmail(c *fiber.Ctx) error {
	input, err := utils.GetValidatedBody[requests.VerifyEmailInput](c)
	if err != nil {
		return err
	}
	token := input.Token
	if token == "" {
		return utils.ErrorResponse(c, utils.ErrorResponseFormat{
			Code:    fiber.StatusUnauthorized,
			Message: "Missing token",
		})
	}

	claims, err := utils.ValidateToken(token)
	if err != nil {
		return utils.ErrorResponse(c, utils.ErrorResponseFormat{
			Code:       fiber.StatusUnauthorized,
			Message:    "Invalid token",
			DevMessage: err.Error(),
			Details:    err,
		})
	}

	user := models.User{}
	if err := ctrl.DB.Where("email = ?", claims["email"]).First(&user).Error; err != nil {
		return utils.ErrorResponse(c, utils.ErrorResponseFormat{
			Code:       fiber.StatusNotFound,
			Message:    "User not found",
			DevMessage: err.Error(),
			Details:    err,
		})
	}
	if user.EmailVerifiedAt != nil {
		return utils.SuccessResponse(c, utils.SuccessResponseFormat{
			Code:    fiber.StatusOK,
			Message: "Email already verified",
		})
	}
	emailVerifiedAt := time.Now()
	user.EmailVerifiedAt = &emailVerifiedAt
	if result := ctrl.DB.Save(&user); result.Error != nil {
		return utils.ErrorResponse(c, utils.ErrorResponseFormat{
			Code:       fiber.StatusInternalServerError,
			Message:    "Failed to save user",
			DevMessage: result.Error.Error(),
			Details:    result.Error,
		})
	}
	return utils.SuccessResponse(c, utils.SuccessResponseFormat{
		Code:    fiber.StatusOK,
		Message: "Email verified successfully",
	})
}
