package middleware

import (
	"errors"
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type ErrorResponse struct {
	Success bool     `json:"success"`
	Message string   `json:"message"`
	Errors  []string `json:"errors,omitempty"`
}

func ErrorHandler() fiber.ErrorHandler {
	return func(c *fiber.Ctx, err error) error {
		code := fiber.StatusInternalServerError
		message := "Internal Server Error"
		var validationErrors []string

		var fiberErr *fiber.Error
		if errors.As(err, &fiberErr) {
			code = fiberErr.Code
			message = fiberErr.Message
		}

		var valErrs validator.ValidationErrors
		if errors.As(err, &valErrs) {
			code = fiber.StatusBadRequest
			message = "Validation failed"
			for _, fieldErr := range valErrs {
				validationErrors = append(validationErrors, fmt.Sprintf("Field '%s' failed validation for tag '%s'", fieldErr.Field(), fieldErr.Tag()))
			}
		}

		if code == fiber.StatusInternalServerError && message == "Internal Server Error" && err != nil {
			message = err.Error()
		}

		return c.Status(code).JSON(ErrorResponse{
			Success: false,
			Message: message,
			Errors:  validationErrors,
		})
	}
}
