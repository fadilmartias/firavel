package utils

import "github.com/gofiber/fiber/v2"

type SuccessResponseFormat struct {
	Code       int
	Message    string
	Data       any
	Pagination *Pagination
	Meta       any
}

type OrderedSuccessResponse struct {
	Success    bool   `json:"success"`
	Message    string `json:"message"`
	Meta       any    `json:"meta,omitempty"`
	Pagination any    `json:"pagination,omitempty"`
	Data       any    `json:"data,omitempty"`
}

type ErrorResponseFormat struct {
	Code    int
	Message string
	Details any
}

type OrderedErrorResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Details any    `json:"details,omitempty"`
}

// SuccessResponse mengirim response JSON standar untuk sukses
func SuccessResponse(c *fiber.Ctx, params SuccessResponseFormat) error {
	response := OrderedSuccessResponse{
		Success:    true,
		Message:    params.Message,
		Data:       params.Data,
		Pagination: params.Pagination,
		Meta:       params.Meta,
	}
	return c.Status(params.Code).JSON(response)
}

// ErrorResponse mengirim response JSON standar untuk error
func ErrorResponse(c *fiber.Ctx, params ErrorResponseFormat) error {
	response := OrderedErrorResponse{
		Success: false,
		Message: params.Message,
	}
	if params.Details != nil {
		response.Details = params.Details
	}
	return c.Status(params.Code).JSON(response)
}
