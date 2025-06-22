package controllers_v0

import (
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"time"

	"goravel/app/processors"
	"goravel/app/registry"
	"goravel/app/utils" // Ganti dengan path utils Anda

	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type GenericController struct {
	DB    *gorm.DB
	Redis *redis.Client
}

func NewGenericController(db *gorm.DB, redis *redis.Client) *GenericController {
	return &GenericController{DB: db, Redis: redis}
}

// Index menangani GET /:model
func (ctrl *GenericController) Index(c *fiber.Ctx) error {
	modelName := c.Params("model")
	modelInfo, err := registry.GetModel(modelName)
	if err != nil {
		return utils.ErrorResponse(c, utils.ErrorResponseFormat{Code: http.StatusNotFound, Message: err.Error()})
	}

	queryParams := c.Request().URI().QueryArgs()
	urlValues := make(url.Values)
	queryParams.VisitAll(func(key, value []byte) {
		urlValues.Add(string(key), string(value))
	})
	params := utils.NewQueryParams(urlValues)
	tx := ctrl.DB.Model(modelInfo.Instance)
	tx = utils.BuildGormQuery(tx, urlValues, false)

	cacheKey := fmt.Sprintf("api:%s:%s", modelName, c.Request().URI().QueryString())
	apiResponse, err := utils.FetchAndCacheDynamic(
		c.UserContext(), ctrl.Redis, tx, params, cacheKey, 1*time.Minute, false,
		modelInfo.Instance, modelInfo.NewSlice,
	)

	if err != nil {
		return utils.ErrorResponse(c, utils.ErrorResponseFormat{Code: http.StatusInternalServerError, Message: err.Error()})
	}

	switch resp := apiResponse.(type) {
	case utils.PaginatedResponse[any]:
		return utils.SuccessResponse(c, utils.SuccessResponseFormat{
			Code: http.StatusOK, Message: "Data retrieved successfully", Data: resp.Data, Pagination: &resp.Pagination,
		})
	default:
		return utils.ErrorResponse(c, utils.ErrorResponseFormat{Code: http.StatusInternalServerError, Message: "Invalid response type"})
	}
}

// Show menangani GET /:model/:id
func (ctrl *GenericController) Show(c *fiber.Ctx) error {
	modelName := c.Params("model")
	id := c.Params("id")
	modelInfo, err := registry.GetModel(modelName)
	if err != nil {
		return utils.ErrorResponse(c, utils.ErrorResponseFormat{Code: http.StatusNotFound, Message: err.Error()})
	}

	urlValues := make(url.Values)
	queryArgs := c.Request().URI().QueryArgs()
	queryArgs.VisitAll(func(key, value []byte) {
		urlValues.Add(string(key), string(value))
	})
	params := utils.NewQueryParams(urlValues)

	// Tambahkan filter by ID atau slug
	tx := ctrl.DB.Model(modelInfo.Instance)
	if c.Query("slug") == "true" {
		tx = tx.Where("slug = ?", id)
	} else {
		tx = tx.Where("id = ?", id)
	}
	tx = utils.BuildGormQuery(tx, urlValues, true)

	cacheKey := fmt.Sprintf("api:%s:%s", modelName, id)
	apiResponse, err := utils.FetchAndCacheDynamic(
		c.UserContext(), ctrl.Redis, tx, params, cacheKey, 5*time.Minute, true,
		modelInfo.Instance, modelInfo.NewSlice,
	)

	if err != nil {
		return utils.ErrorResponse(c, utils.ErrorResponseFormat{Code: http.StatusInternalServerError, Message: err.Error()})
	}

	switch resp := apiResponse.(type) {
	case utils.SingleResponse[any]:
		baseURL := fmt.Sprintf("%s://%s", c.Protocol(), c.Hostname())
		metaData, _ := processors.GenericPostProcessor(resp.Data, baseURL)

		return utils.SuccessResponse(c, utils.SuccessResponseFormat{
			Code:    http.StatusOK,
			Message: "Data retrieved successfully",
			Data:    resp.Data,
			Meta:    metaData, // Kirim meta yang sudah diproses
		})
	default:
		return utils.ErrorResponse(c, utils.ErrorResponseFormat{Code: http.StatusInternalServerError, Message: "Invalid response type"})
	}
}

func (ctrl *GenericController) Store(c *fiber.Ctx) error {
	// 1. Dapatkan info model dari registry
	modelName := c.Params("model")
	modelInfo, err := registry.GetModel(modelName)
	if err != nil {
		return utils.ErrorResponse(c, utils.ErrorResponseFormat{Code: http.StatusNotFound, Message: err.Error()})
	}

	// 2. Buat instance baru dari model menggunakan reflection
	// `newInstance` sekarang adalah pointer ke struct kosong, misal: &models.User{}
	newInstance := reflect.New(reflect.TypeOf(modelInfo.Instance).Elem()).Interface()

	// 3. Parse request body JSON ke dalam instance baru tersebut
	if err := c.BodyParser(newInstance); err != nil {
		return utils.ErrorResponse(c, utils.ErrorResponseFormat{Code: http.StatusBadRequest, Message: "Invalid request body"})
	}

	// TODO: Implementasikan validasi data di sini menggunakan library seperti `go-playground/validator`

	// Placeholder untuk logika file upload
	// utils.HandleFileUpload(c, modelName)

	// 4. Set ID jika belum ada (sesuai logika JS Anda)
	// Kita perlu reflection lagi untuk mengakses dan set field 'Id'
	val := reflect.ValueOf(newInstance).Elem()
	idField := val.FieldByName("Id")
	if idField.IsValid() && idField.Kind() == reflect.String && idField.String() == "" {
		idField.SetString(utils.GenerateShortID(7))
	}

	// 5. Simpan instance ke database
	if result := ctrl.DB.Create(newInstance); result.Error != nil {
		return utils.ErrorResponse(c, utils.ErrorResponseFormat{Code: http.StatusInternalServerError, Message: result.Error.Error()})
	}

	// 6. Kembalikan data yang baru dibuat
	return utils.SuccessResponse(c, utils.SuccessResponseFormat{
		Code:    http.StatusCreated,
		Message: "Data created successfully",
		Data:    newInstance,
	})
}

// Update menangani PUT /:model/:id
func (ctrl *GenericController) Update(c *fiber.Ctx) error {
	// 1. Dapatkan info model dan parameter
	modelName := c.Params("model")
	id := c.Params("id")
	modelInfo, err := registry.GetModel(modelName)
	if err != nil {
		return utils.ErrorResponse(c, utils.ErrorResponseFormat{Code: http.StatusNotFound, Message: err.Error()})
	}

	// 2. Buat instance untuk menampung data update dari body
	updateData := reflect.New(reflect.TypeOf(modelInfo.Instance).Elem()).Interface()
	if err := c.BodyParser(updateData); err != nil {
		return utils.ErrorResponse(c, utils.ErrorResponseFormat{Code: http.StatusBadRequest, Message: "Invalid request body"})
	}

	// 3. Lakukan update di database
	// GORM secara otomatis hanya akan mengupdate field yang tidak "zero-value" dari `updateData`
	// Ini sangat efisien dan aman.
	result := ctrl.DB.Model(modelInfo.Instance).Where("id = ?", id).Updates(updateData)
	if result.Error != nil {
		return utils.ErrorResponse(c, utils.ErrorResponseFormat{Code: http.StatusInternalServerError, Message: result.Error.Error()})
	}

	// 4. Periksa apakah ada baris yang terpengaruh
	if result.RowsAffected == 0 {
		return utils.ErrorResponse(c, utils.ErrorResponseFormat{Code: http.StatusNotFound, Message: "Data not found or no changes made"})
	}

	return utils.SuccessResponse(c, utils.SuccessResponseFormat{
		Code:    http.StatusOK,
		Message: "Data updated successfully",
		// Mengembalikan jumlah baris yang diupdate, sesuai perilaku GORM
		Data: fiber.Map{"rows_affected": result.RowsAffected},
	})
}

// Destroy menangani DELETE /:model/:id
func (ctrl *GenericController) Destroy(c *fiber.Ctx) error {
	// 1. Dapatkan info model dan parameter
	modelName := c.Params("model")
	id := c.Params("id")
	modelInfo, err := registry.GetModel(modelName)
	if err != nil {
		return utils.ErrorResponse(c, utils.ErrorResponseFormat{Code: http.StatusNotFound, Message: err.Error()})
	}

	// 2. Hapus data dari database
	// GORM akan melakukan soft delete jika model memiliki field `DeletedAt`
	result := ctrl.DB.Where("id = ?", id).Delete(modelInfo.Instance)
	if result.Error != nil {
		return utils.ErrorResponse(c, utils.ErrorResponseFormat{Code: http.StatusInternalServerError, Message: result.Error.Error()})
	}

	// 3. Periksa apakah ada baris yang terpengaruh
	if result.RowsAffected == 0 {
		return utils.ErrorResponse(c, utils.ErrorResponseFormat{Code: http.StatusNotFound, Message: "Data not found"})
	}

	return utils.SuccessResponse(c, utils.SuccessResponseFormat{
		Code:    http.StatusOK,
		Message: "Data deleted successfully",
	})
}
