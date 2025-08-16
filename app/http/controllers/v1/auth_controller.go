package controllers_v1

import (
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/fadilmartias/firavel/app/jobs"
	"github.com/fadilmartias/firavel/app/models"
	"github.com/fadilmartias/firavel/app/requests"
	"github.com/fadilmartias/firavel/app/utils"
	"github.com/fadilmartias/firavel/config"
	"github.com/go-resty/resty/v2"
	"github.com/golang-jwt/jwt"
	"github.com/tidwall/gjson"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm" // TAMBAHKAN IMPORT INI
)

type AuthController struct {
	BaseController
	DB    *gorm.DB // Tambahkan ini untuk menyimpan koneksi DB
	Redis *config.RedisClient
}

// Ubah fungsi NewAuthController untuk menerima koneksi DB
func NewAuthController(db *gorm.DB, redis *config.RedisClient) *AuthController {
	return &AuthController{DB: db, Redis: redis}
}

func (ctrl *AuthController) Me(c *fiber.Ctx) error {
	user, ok := c.Locals("user").(jwt.MapClaims)
	if !ok {
		return utils.ErrorResponse(c, utils.ErrorResponseFormat{
			Code:    fiber.StatusUnauthorized,
			Message: "Invalid token",
		})
	}
	return utils.SuccessResponse(c, utils.SuccessResponseFormat{
		Code:    fiber.StatusOK,
		Message: "User retrieved successfully",
		Data:    user,
	})
}

func (ctrl *AuthController) Logout(c *fiber.Ctx) error {
	user, ok := c.Locals("user").(jwt.MapClaims)
	if !ok {
		return utils.ErrorResponse(c, utils.ErrorResponseFormat{
			Code:    fiber.StatusUnauthorized,
			Message: "Invalid token",
		})
	}

	var userDB models.User
	if err := ctrl.DB.Where("email = ?", user["email"]).First(&userDB).Error; err != nil {
		return utils.ErrorResponse(c, utils.ErrorResponseFormat{
			Code:    fiber.StatusNotFound,
			Message: "User not found",
			Details: nil,
		})
	}

	userDB.RefreshToken = nil

	c.Cookie(&fiber.Cookie{
		Name:     "refresh_token_" + os.Getenv("APP_ENV"),
		Value:    "",
		Expires:  time.Now().Add(-time.Hour * 24),
		HTTPOnly: true,
		Domain:   os.Getenv("FE_DOMAIN"),
		SameSite: fiber.CookieSameSiteNoneMode,
		Secure:   true,
	})
	c.Cookie(&fiber.Cookie{
		Name:     "access_token_" + os.Getenv("APP_ENV"),
		Value:    "",
		Expires:  time.Now().Add(-time.Hour * 1),
		HTTPOnly: true,
		Domain:   os.Getenv("FE_DOMAIN"),
		SameSite: fiber.CookieSameSiteNoneMode,
		Secure:   true,
	})

	if err := ctrl.DB.Save(&userDB).Error; err != nil {
		return utils.ErrorResponse(c, utils.ErrorResponseFormat{
			Code:    fiber.StatusInternalServerError,
			Message: "Failed to logout user",
			Details: nil,
		})
	}

	return utils.SuccessResponse(c, utils.SuccessResponseFormat{
		Code:    fiber.StatusOK,
		Message: "Logged out successfully",
	})
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
	User         userResponse `json:"user"`
	AccessToken  string       `json:"access_token"`
	RefreshToken string       `json:"refresh_token"`
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
			Code:       fiber.StatusNotFound,
			Message:    "Pengguna tidak ditemukan",
			DevMessage: err.Error(),
			Details:    err,
		})
	}

	// Bandingkan password input dengan hash di DB
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		return utils.ErrorResponse(c, utils.ErrorResponseFormat{
			Code:       fiber.StatusBadRequest,
			Message:    "Email atau password salah",
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
			Message:    "Gagal menghasilkan access token",
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
			Message:    "Gagal menghasilkan refresh token",
			DevMessage: err.Error(),
			Details:    err,
		})
	}

	accessCookie := fiber.Cookie{
		Name:     "access_token_" + os.Getenv("APP_ENV"),
		Value:    accessToken,
		Expires:  time.Now().Add(time.Hour * 1),
		HTTPOnly: true,
		Domain:   os.Getenv("FE_DOMAIN"),
		SameSite: fiber.CookieSameSiteNoneMode,
		Secure:   true,
	}
	c.Cookie(&accessCookie)

	refreshCookie := fiber.Cookie{
		Name:     "refresh_token_" + os.Getenv("APP_ENV"),
		Value:    refreshToken,
		Expires:  time.Now().Add(time.Hour * 24),
		HTTPOnly: true,
		Domain:   os.Getenv("FE_DOMAIN"),
		SameSite: fiber.CookieSameSiteNoneMode,
		Secure:   true,
	}
	c.Cookie(&refreshCookie)

	user.RefreshToken = &refreshToken
	if err := ctrl.DB.Save(&user).Error; err != nil {
		return utils.ErrorResponse(c, utils.ErrorResponseFormat{
			Code:       fiber.StatusInternalServerError,
			Message:    "Gagal menyimpan refresh token",
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
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
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
			Message:    "Pengguna tidak ditemukan",
			DevMessage: err.Error(),
			Details:    err,
		})
	}

	if err := ctrl.DB.Where("email = ?", input.Email).Where("expired_at > ?", time.Now()).First(&passwordResetToken).Error; err == nil {
		return utils.SuccessResponse(c, utils.SuccessResponseFormat{
			Code:    fiber.StatusOK,
			Message: "Token reset password sudah dikirim",
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
	task, err := jobs.NewEmailResetPasswordTask(user.Email, token)
	if err != nil {
		log.Fatal("failed to create task:", err)
	}
	info, err := jobs.AsynqClient.Enqueue(task)
	if err != nil {
		log.Fatal("failed to enqueue task:", err)
	}
	fmt.Println("enqueued task:", info)
	return utils.SuccessResponse(c, utils.SuccessResponseFormat{
		Code:    fiber.StatusOK,
		Message: "Token reset password berhasil dikirim",
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
			Message:    "Pengguna tidak ditemukan",
			DevMessage: err.Error(),
			Details:    err,
		})
	}

	user.HashPassword(input.Password)
	if err := ctrl.DB.Save(&user).Error; err != nil {
		return utils.ErrorResponse(c, utils.ErrorResponseFormat{
			Code:       fiber.StatusInternalServerError,
			Message:    "Gagal menyimpan user",
			DevMessage: err.Error(),
			Details:    err,
		})
	}

	now := time.Now()

	passwordResetToken.UsedAt = &now
	ctrl.DB.Save(&passwordResetToken)

	return utils.SuccessResponse(c, utils.SuccessResponseFormat{
		Code:    fiber.StatusOK,
		Message: "Password reset berhasil",
	})
}

func (ctrl *AuthController) SendEmailVerification(c *fiber.Ctx) error {
	id := c.Locals("user").(jwt.MapClaims)["id"].(string)
	var user models.User
	if err := ctrl.DB.Where("id = ?", id).First(&user).Error; err != nil {
		return utils.ErrorResponse(c, utils.ErrorResponseFormat{
			Code:    fiber.StatusNotFound,
			Message: "Pengguna tidak ditemukan",
			Details: nil,
		})
	}
	if user.EmailVerifiedAt != nil {
		return utils.SuccessResponse(c, utils.SuccessResponseFormat{
			Code:    fiber.StatusOK,
			Message: "Email sudah terverifikasi",
		})
	}
	_, err := ctrl.Redis.Get(c.UserContext(), fmt.Sprintf("email_verification_token:%s", user.Email))
	if err == nil {
		return utils.ErrorResponse(c, utils.ErrorResponseFormat{
			Code:    fiber.StatusConflict,
			Message: "Email verifikasi sudah dikirim",
		})
	}

	jwtToken, _ := utils.GenerateToken(map[string]any{
		"email": user.Email,
	}, time.Minute*5)

	if err := ctrl.Redis.Set(c.UserContext(), fmt.Sprintf("email_verification_token:%s", user.Email), jwtToken, time.Minute*5); err != nil {
		return utils.ErrorResponse(c, utils.ErrorResponseFormat{
			Code:       fiber.StatusInternalServerError,
			Message:    "Gagal menyimpan token verifikasi email",
			DevMessage: err.Error(),
			Details:    err,
		})
	}

	// Send email
	task, err := jobs.NewEmailVerificationTask(user.Email, jwtToken)
	if err != nil {
		log.Fatal("failed to create task:", err)
	}
	info, err := jobs.AsynqClient.Enqueue(task)
	if err != nil {
		log.Fatal("failed to enqueue task:", err)
	}
	fmt.Println("enqueued task:", info)

	return utils.SuccessResponse(c, utils.SuccessResponseFormat{
		Code:    fiber.StatusOK,
		Message: "Email verifikasi berhasil dikirim",
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
			Message: "Token tidak valid",
		})
	}

	claims, err := utils.ValidateToken(token)
	if err != nil {
		return utils.ErrorResponse(c, utils.ErrorResponseFormat{
			Code:       fiber.StatusUnauthorized,
			Message:    "Token tidak valid",
			DevMessage: err.Error(),
			Details:    err,
		})
	}

	user := models.User{}
	if err := ctrl.DB.Where("email = ?", claims["email"]).First(&user).Error; err != nil {
		return utils.ErrorResponse(c, utils.ErrorResponseFormat{
			Code:       fiber.StatusNotFound,
			Message:    "Pengguna tidak ditemukan",
			DevMessage: err.Error(),
			Details:    err,
		})
	}
	if user.EmailVerifiedAt != nil {
		return utils.SuccessResponse(c, utils.SuccessResponseFormat{
			Code:    fiber.StatusOK,
			Message: "Email sudah terverifikasi",
		})
	}
	emailVerifiedAt := time.Now()
	user.EmailVerifiedAt = &emailVerifiedAt
	if result := ctrl.DB.Save(&user); result.Error != nil {
		return utils.ErrorResponse(c, utils.ErrorResponseFormat{
			Code:       fiber.StatusInternalServerError,
			Message:    "Gagal menyimpan user",
			DevMessage: result.Error.Error(),
			Details:    result.Error,
		})
	}
	return utils.SuccessResponse(c, utils.SuccessResponseFormat{
		Code:    fiber.StatusOK,
		Message: "Email verified successfully",
	})

}

// RefreshToken mengembalikan token yang baru
func (ctrl *AuthController) RefreshAccessToken(c *fiber.Ctx) error {
	token := c.Cookies("refreshToken")
	if token == "" {
		return utils.ErrorResponse(c, utils.ErrorResponseFormat{
			Code:    fiber.StatusBadRequest,
			Message: "Token tidak ditemukan",
		})
	}
	claims, err := utils.ValidateToken(token)
	if err != nil {
		return utils.ErrorResponse(c, utils.ErrorResponseFormat{
			Code:       fiber.StatusBadRequest,
			Message:    "Token tidak valid",
			DevMessage: err.Error(),
			Details:    err,
		})
	}

	user := models.User{}
	if err := ctrl.DB.Where("email = ?", claims["email"]).First(&user).Error; err != nil {
		return utils.ErrorResponse(c, utils.ErrorResponseFormat{
			Code:       fiber.StatusNotFound,
			Message:    "Pengguna tidak ditemukan",
			DevMessage: err.Error(),
			Details:    err,
		})
	}

	accessToken, err := utils.GenerateToken(map[string]any{
		"id":    user.ID,
		"name":  user.Name,
		"email": user.Email,
		"role":  user.Role,
	}, time.Hour*1)

	refreshToken, err := utils.GenerateToken(map[string]any{
		"id":    user.ID,
		"name":  user.Name,
		"email": user.Email,
		"role":  user.Role,
	}, time.Hour*24)

	if err := ctrl.DB.Model(&user).Update("refresh_token", refreshToken).Error; err != nil {
		return utils.ErrorResponse(c, utils.ErrorResponseFormat{
			Code:       fiber.StatusInternalServerError,
			Message:    "Gagal menyimpan refresh token",
			DevMessage: err.Error(),
			Details:    err,
		})
	}
	accessCookie := fiber.Cookie{
		Name:     "access_token_" + os.Getenv("APP_ENV"),
		Value:    accessToken,
		Expires:  time.Now().Add(time.Hour * 1),
		HTTPOnly: true,
		Domain:   os.Getenv("FE_DOMAIN"),
		SameSite: fiber.CookieSameSiteNoneMode,
		Secure:   true,
	}
	refreshCookie := fiber.Cookie{
		Name:     "refresh_token_" + os.Getenv("APP_ENV"),
		Value:    refreshToken,
		Expires:  time.Now().Add(time.Hour * 24),
		HTTPOnly: true,
		Domain:   os.Getenv("FE_DOMAIN"),
		SameSite: fiber.CookieSameSiteNoneMode,
		Secure:   true,
	}
	c.Cookie(&accessCookie)
	c.Cookie(&refreshCookie)

	if err != nil {
		return utils.ErrorResponse(c, utils.ErrorResponseFormat{
			Code:       fiber.StatusInternalServerError,
			Message:    "Gagal menghasilkan token",
			DevMessage: err.Error(),
			Details:    err,
		})
	}

	return utils.SuccessResponse(c, utils.SuccessResponseFormat{
		Code:    fiber.StatusOK,
		Message: "Token refreshed successfully",
		Data: fiber.Map{
			"access_token":  accessToken,
			"refresh_token": refreshToken,
		},
	})
}

func (ctrl *AuthController) GoogleRedirect(c *fiber.Ctx) error {
	clientID := os.Getenv("GOOGLE_CLIENT_ID")
	redirectURI := os.Getenv("APP_URL") + "/v1/auth/google/callback"
	authURL := fmt.Sprintf(
		"https://accounts.google.com/o/oauth2/v2/auth?client_id=%s&redirect_uri=%s&response_type=code&scope=email profile",
		clientID, redirectURI,
	)
	return c.Redirect(authURL)
}

func (ctrl *AuthController) GoogleCallback(c *fiber.Ctx) error {
	code := c.Query("code")
	if code == "" {
		return c.Status(400).SendString("Code not found")
	}

	clientID := os.Getenv("GOOGLE_CLIENT_ID")
	clientSecret := os.Getenv("GOOGLE_CLIENT_SECRET")
	redirectURI := os.Getenv("APP_URL") + "/v1/auth/google/callback"

	// Tukar code ke token
	resp, err := resty.New().
		R().
		SetFormData(map[string]string{
			"code":          code,
			"client_id":     clientID,
			"client_secret": clientSecret,
			"redirect_uri":  redirectURI,
			"grant_type":    "authorization_code",
		}).
		Post("https://oauth2.googleapis.com/token")
	if err != nil {
		return err
	}

	accessTokenGoogle := gjson.GetBytes(resp.Body(), "access_token").String()

	// Ambil profile
	userResp, err := resty.New().
		R().
		SetAuthToken(accessTokenGoogle).
		Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		return err
	}

	email := gjson.GetBytes(userResp.Body(), "email").String()
	name := gjson.GetBytes(userResp.Body(), "name").String()

	user := models.User{}
	if err := ctrl.DB.Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Buat user baru
			emailVerifiedAt := time.Now()
			user = models.User{
				Name:            name,
				Email:           email,
				Phone:           "",
				Role:            "user",
				EmailVerifiedAt: &emailVerifiedAt,
			}
			if err := ctrl.DB.Create(&user).Error; err != nil {
				return utils.ErrorResponse(c, utils.ErrorResponseFormat{
					Code:       fiber.StatusInternalServerError,
					Message:    "Gagal membuat user",
					Details:    err,
					DevMessage: err.Error(),
				})
			}
		} else {
			return utils.ErrorResponse(c, utils.ErrorResponseFormat{
				Code:       fiber.StatusInternalServerError,
				Message:    "Gagal memeriksa user",
				Details:    err,
				DevMessage: err.Error(),
			})
		}
	} else {
		// Update user
		if err := ctrl.DB.Model(&user).Updates(models.User{
			Name: name,
		}).Error; err != nil {
			return utils.ErrorResponse(c, utils.ErrorResponseFormat{
				Code:       fiber.StatusInternalServerError,
				Message:    "Gagal memperbarui user",
				Details:    err,
				DevMessage: err.Error(),
			})
		}
	}

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
			Message:    "Gagal menghasilkan access token",
			Details:    err,
			DevMessage: err.Error(),
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
			Message:    "Gagal menghasilkan refresh token",
			Details:    err,
			DevMessage: err.Error(),
		})
	}
	user.RefreshToken = &refreshToken
	if err := ctrl.DB.Save(&user).Error; err != nil {
		return utils.ErrorResponse(c, utils.ErrorResponseFormat{
			Code:       fiber.StatusInternalServerError,
			Message:    "Gagal menyimpan refresh token",
			Details:    err,
			DevMessage: err.Error(),
		})
	}

	accessCookie := fiber.Cookie{
		Name:     "access_token_" + os.Getenv("APP_ENV"),
		Value:    accessToken,
		Expires:  time.Now().Add(time.Hour * 1),
		HTTPOnly: true,
		Domain:   os.Getenv("FE_DOMAIN"),
		SameSite: fiber.CookieSameSiteNoneMode,
		Secure:   true,
	}
	c.Cookie(&accessCookie)

	refreshCookie := fiber.Cookie{
		Name:     "refresh_token_" + os.Getenv("APP_ENV"),
		Value:    refreshToken,
		Expires:  time.Now().Add(time.Hour * 24),
		HTTPOnly: true,
		Domain:   os.Getenv("FE_DOMAIN"),
		SameSite: fiber.CookieSameSiteNoneMode,
		Secure:   true,
	}
	c.Cookie(&refreshCookie)

	if user.Role == "admin" {
		return c.Redirect(os.Getenv("FE_URL") + "/admin/dashboard")
	}
	// Redirect ke FE
	return c.Redirect(os.Getenv("FE_URL"))
}
