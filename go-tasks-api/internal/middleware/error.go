package middleware

import (
	"log"

	"github.com/Kathecam/go-tasks-api/internal/dto"
	"github.com/Kathecam/go-tasks-api/internal/errors"
	"github.com/gofiber/fiber/v2"
)

// ErrorHandler middleware para manejo centralizado de errores
func ErrorHandler() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Ejecutar el handler siguiente
		err := c.Next()

		if err == nil {
			return nil
		}

		// Verificar si es nuestro AppError personalizado
		if appErr, ok := errors.IsAppError(err); ok {
			return c.Status(appErr.Code).JSON(dto.ErrorResponse{
				Error:   "Application Error",
				Message: appErr.Message,
				Code:    appErr.Code,
				Details: appErr.Details,
			})
		}

		// Verificar si es un error de Fiber
		if fiberErr, ok := err.(*fiber.Error); ok {
			return c.Status(fiberErr.Code).JSON(dto.ErrorResponse{
				Error:   "HTTP Error",
				Message: fiberErr.Message,
				Code:    fiberErr.Code,
			})
		}

		// Error desconocido - log para debugging
		log.Printf("Unhandled error: %v", err)

		// Retornar error genérico 500
		return c.Status(500).JSON(dto.ErrorResponse{
			Error:   "Internal Server Error",
			Message: "An unexpected error occurred",
			Code:    500,
		})
	}
}
