package controllers_v0

import (
	"fmt"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/fadilmartias/firavel/app/utils" // Ganti dengan path utils Anda
	"github.com/fadilmartias/firavel/config"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type UploadController struct {
	DB    *gorm.DB
	Redis *config.RedisClient
}

func NewUploadController(db *gorm.DB, redis *config.RedisClient) *UploadController {
	return &UploadController{DB: db, Redis: redis}
}

func (ctrl *UploadController) Upload(c *fiber.Ctx) error {
	form, err := c.MultipartForm()
	if err != nil {
		return utils.ErrorResponse(c, utils.ErrorResponseFormat{
			Code:    fiber.StatusBadRequest,
			Message: "Failed to parse form",
		})
	}

	files := form.File
	acceptedFields := []string{"image", "file", "tinymce"}

	var uploadedFile *multipart.FileHeader
	var fieldName string

	for _, field := range acceptedFields {
		fileHeaders, ok := files[field]
		if ok && len(fileHeaders) > 0 {
			uploadedFile = fileHeaders[0]
			fieldName = field
			break
		}
	}

	if uploadedFile == nil {
		return utils.ErrorResponse(c, utils.ErrorResponseFormat{
			Code:    fiber.StatusBadRequest,
			Message: "No valid file field found (image, file)",
		})
	}

	// Validasi Content-Type
	contentType := uploadedFile.Header.Get("Content-Type")
	allowedTypes := map[string]string{
		"image/jpeg":    ".jpg",
		"image/png":     ".png",
		"image/webp":    ".webp",
		"image/gif":     ".gif",
		"image/svg+xml": ".svg",
	}
	ext := filepath.Ext(uploadedFile.Filename)
	if ext == "" || ext == "blob" {
		var ok bool
		ext, ok = allowedTypes[contentType]
		if !ok {
			return utils.ErrorResponse(c, utils.ErrorResponseFormat{Code: fiber.StatusBadRequest, Message: "Unsupported file type"})
		}
	}

	// Buat folder public/tmp/images jika belum ada
	var uploadDir string
	if fieldName == "tinymce" {
		uploadDir = "public/uploads/tinymce"
	} else {
		uploadDir = "public/uploads/tmp/" + strings.ToLower(fieldName) + "s"
	}
	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		return utils.ErrorResponse(c, utils.ErrorResponseFormat{
			Code:    fiber.StatusInternalServerError,
			Message: "Failed to create upload directory",
		})
	}

	// Rename file dengan timestamp dan random string
	timestamp := time.Now().Unix()
	randomStr := utils.GenerateShortID(6) // pastikan ini menghasilkan string acak
	newFileName := fmt.Sprintf("%s-%d-%s%s", fieldName, timestamp, randomStr, ext)
	savePath := filepath.Join(uploadDir, newFileName)

	// Simpan file
	if err := c.SaveFile(uploadedFile, savePath); err != nil {
		return utils.ErrorResponse(c, utils.ErrorResponseFormat{
			Code:    fiber.StatusInternalServerError,
			Message: "Failed to save file",
		})
	}

	// File URL
	var fileURL string
	if fieldName == "tinymce" {
		fileURL = fmt.Sprintf(os.Getenv("APP_URL")+"/uploads/tinymce/%s", newFileName)
	} else {
		fileURL = fmt.Sprintf(os.Getenv("APP_URL")+"/uploads/tmp/%ss/%s", strings.ToLower(fieldName), newFileName)
	}

	// Respons sukses
	return utils.SuccessResponse(c, utils.SuccessResponseFormat{
		Code:    fiber.StatusOK,
		Message: "Upload successfully",
		Data: fiber.Map{
			"file_url":  fileURL,
			"file_name": newFileName,
		},
	})
}
