package utils

import "github.com/gofiber/fiber/v2"

type SuccessResponseFormat struct {
	Code       int
	Message    string
	Data       any
	Pagination *Pagination
}

type ErrorResponseFormat struct {
	Code    int
	Message string
	Details any
}

// SuccessResponse mengirim response JSON standar untuk sukses
func SuccessResponse(c *fiber.Ctx, params SuccessResponseFormat) error {
	response := fiber.Map{
		"success": true,
		"message": params.Message,
		"data":    params.Data,
	}
	if params.Pagination != nil {
		response["pagination"] = params.Pagination
	}
	return c.Status(params.Code).JSON(response)
}

// ErrorResponse mengirim response JSON standar untuk error
func ErrorResponse(c *fiber.Ctx, params ErrorResponseFormat) error {
	response := fiber.Map{
		"success": false,
		"message": params.Message,
	}
	if params.Details != nil {
		response["details"] = params.Details
	}
	return c.Status(params.Code).JSON(response)
}
