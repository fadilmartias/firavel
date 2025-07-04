package utils

import (
	"os"
	"runtime/debug"

	"github.com/gofiber/fiber/v2"
)

type SuccessResponseFormat struct {
	Code       int
	Message    string
	Data       any
	Pagination *Pagination
	Meta       any
}

type OrderedSuccessResponse struct {
	Success    bool        `json:"success"`
	Message    string      `json:"message"`
	Meta       any         `json:"meta,omitempty"`
	Pagination *Pagination `json:"pagination,omitempty"`
	Data       any         `json:"data,omitempty"`
}

type ErrorResponseFormat struct {
	Code       int
	Message    string
	DevMessage string
	Details    any
	Trace      string
}

type OrderedErrorResponse struct {
	Success    bool   `json:"success"`
	Message    string `json:"message"`
	DevMessage string `json:"dev_message,omitempty"`
	Details    any    `json:"details,omitempty"`
	Trace      string `json:"trace,omitempty"`
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
	if os.Getenv("APP_ENV") != "production" && params.DevMessage != "" {
		response.DevMessage = params.DevMessage
	}
	if os.Getenv("APP_ENV") != "production" {
		response.Trace = string(debug.Stack())
	}
	return c.Status(params.Code).JSON(response)
}
