package utils

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

func HandleFileUpload(c *fiber.Ctx, modelName string, field string, typeFile string) (string, error) {
	// TODO: Implementasikan logika upload file di sini.
	file, err := c.FormFile(field)
	if err != nil {
		return "", err
	}

	// Simpan file ke folder public/uploads
	filename := fmt.Sprintf("%s_%s", modelName, file.Filename)
	if err := c.SaveFile(file, fmt.Sprintf("./public/uploads/%s/%s", typeFile, filename)); err != nil {
		return "", err
	}

	return fmt.Sprintf("/uploads/%s/%s", typeFile, filename), nil
}
