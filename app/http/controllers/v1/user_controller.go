package controllers_v1

import (
	"goravel/app/models"
	"goravel/app/utils"
	"net/http"
	"net/url"
	"time"

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

// Ganti `models.User` dengan struct model yang ingin Anda ambil
func (ctrl *UserController) Flex(c *fiber.Ctx) error {
	// --- LANGKAH 1: Parse URL Params ---
	queryParams := c.Request().URI().QueryArgs()
	urlValues := make(url.Values)
	queryParams.VisitAll(func(key, value []byte) {
		urlValues.Add(string(key), string(value))
	})
	params := utils.NewQueryParams(urlValues) // Buat struct params

	// --- LANGKAH 2: Bangun Kueri GORM ---
	// Tentukan model yang akan di-query
	tx := ctrl.DB.Model(&models.User{})
	// Bangun kueri dinamis dari params
	tx = utils.BuildGormQuery(tx, urlValues, false)

	// --- LANGKAH 3: Fetch Data dengan Caching ---
	// Buat cache key yang unik berdasarkan query
	cacheKey := "users:" + string(c.Request().URI().QueryString())
	cacheDuration := 5 * time.Minute // Contoh durasi cache

	// Panggil FetchAndCache dengan argumen dan tipe generic yang benar
	response, err := utils.FetchAndCache[models.User](
		c.UserContext(), // 1. Context dari Fiber
		ctrl.Redis,      // 2. Klien Redis
		tx,              // 3. Objek kueri GORM yang sudah dibangun
		params,          // 4. Struct QueryParams
		cacheKey,        // 5. Kunci Cache
		cacheDuration,   // 6. Durasi Cache
		false,           // 7. isSingle
	)

	if err != nil {
		return utils.ErrorResponse(c, utils.ErrorResponseFormat{
			Code:    http.StatusInternalServerError,
			Message: "Gagal mengambil data",
			Details: err.Error(),
		})
	}

	switch resp := response.(type) {

	// Kasus jika respons memiliki paginasi
	case utils.PaginatedResponse[models.User]:
		return utils.SuccessResponse(c, utils.SuccessResponseFormat{
			Code:       http.StatusOK,
			Message:    "Data berhasil diambil",
			Data:       resp.Data,        // <-- Ambil data dari dalam struct
			Pagination: &resp.Pagination, // <-- Ambil pagination dari dalam struct
		})

	// Kasus jika respons adalah data tunggal
	case utils.SingleResponse[models.User]:
		return utils.SuccessResponse(c, utils.SuccessResponseFormat{
			Code:       http.StatusOK,
			Message:    "Data berhasil diambil",
			Data:       resp.Data, // <-- Ambil data dari dalam struct
			Pagination: nil,       // <-- Tidak ada pagination
		})

	// Kasus default jika tipe tidak terduga
	default:
		return utils.ErrorResponse(c, utils.ErrorResponseFormat{
			Code:    http.StatusInternalServerError,
			Message: "Tipe respons tidak valid",
		})
	}
}
