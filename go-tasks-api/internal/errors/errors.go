package errors

import (
	"fmt"
	"net/http"
)

// AppError representa un error de aplicación con contexto
type AppError struct {
	Code    int               `json:"code"`
	Message string            `json:"message"`
	Details map[string]string `json:"details,omitempty"`
	Err     error             `json:"-"` // Error original (no se serializa)
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

// New crea un nuevo AppError
func New(code int, message string) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
	}
}

// Wrap envuelve un error existente con contexto
func Wrap(err error, code int, message string) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}

// WithDetails agrega detalles al error
func (e *AppError) WithDetails(details map[string]string) *AppError {
	e.Details = details
	return e
}

// Errores comunes predefinidos
var (
	// 400 - Bad Request
	ErrInvalidInput = New(http.StatusBadRequest, "Invalid input data")
	ErrInvalidJSON  = New(http.StatusBadRequest, "Invalid JSON format")

	// 404 - Not Found
	ErrTaskNotFound = New(http.StatusNotFound, "Task not found")
	ErrNotFound     = New(http.StatusNotFound, "Resource not found")

	// 422 - Unprocessable Entity
	ErrValidationFailed = New(http.StatusUnprocessableEntity, "Validation failed")

	// 500 - Internal Server Error
	ErrInternalServer = New(http.StatusInternalServerError, "Internal server error")
	ErrDatabaseError  = New(http.StatusInternalServerError, "Database error")
)

// IsAppError verifica si un error es de tipo AppError
func IsAppError(err error) (*AppError, bool) {
	appErr, ok := err.(*AppError)
	return appErr, ok
}
