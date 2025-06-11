package controllers_v0

import (
	"fmt"
	"net/http"
	"net/url"
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
