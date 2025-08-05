package file_template

const ControllerTemplate = `package controllers_v1

import (
	"github.com/fadilmartias/firavel/app/utils"
	"github.com/fadilmartias/firavel/config"
	
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type {{.Name}} struct {
	BaseController
	DB    *gorm.DB
	Redis *config.RedisClient
}

func New{{.Name}}(db *gorm.DB, redis *config.RedisClient) *{{.Name}} {
	return &{{.Name}}{DB: db, Redis: redis}
}

func (ctrl *{{.Name}}) Index(c *fiber.Ctx) error {
	return utils.SuccessResponse(c, utils.SuccessResponseFormat{
		Message: "Berhasil mendapatkan data {{.LowerName}}",
		Data:    nil,
	})
}

func (ctrl *{{.Name}}) Show(c *fiber.Ctx) error {
	return utils.SuccessResponse(c, utils.SuccessResponseFormat{
		Message: "Berhasil mendapatkan data {{.LowerName}}",
		Data:    nil,
	})
}

func (ctrl *{{.Name}}) Store(c *fiber.Ctx) error {
	return utils.SuccessResponse(c, utils.SuccessResponseFormat{
		Message: "Berhasil menambahkan {{.LowerName}}",
		Data:    nil,
	})
}

func (ctrl *{{.Name}}) Update(c *fiber.Ctx) error {
	return utils.SuccessResponse(c, utils.SuccessResponseFormat{
		Message: "Berhasil mengupdate {{.LowerName}}",
		Data:    nil,
	})
}

func (ctrl *{{.Name}}) Destroy(c *fiber.Ctx) error {
	return utils.SuccessResponse(c, utils.SuccessResponseFormat{
		Message: "Berhasil menghapus {{.LowerName}}",
		Data:    nil,
	})
}

`
