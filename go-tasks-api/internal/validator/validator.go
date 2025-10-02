package validator

import (
	"strings"

	"github.com/Kathecam/go-tasks-api/internal/errors"
	"github.com/go-playground/validator/v10"
)

// Validator encapsula la funcionalidad de validación
type Validator struct {
	validate *validator.Validate
}

// New crea una nueva instancia del validador
func New() *Validator {
	return &Validator{
		validate: validator.New(),
	}
}

// ValidateStruct valida una estructura y retorna errores formateados
func (v *Validator) ValidateStruct(s interface{}) error {
	err := v.validate.Struct(s)
	if err == nil {
		return nil
	}

	// Convertir errores de validación a nuestro formato
	validationErrors := err.(validator.ValidationErrors)
	details := make(map[string]string)

	for _, fieldError := range validationErrors {
		field := strings.ToLower(fieldError.Field())
		message := getValidationMessage(fieldError)
		details[field] = message
	}

	return errors.ErrValidationFailed.WithDetails(details)
}

// getValidationMessage convierte errores de validator a mensajes legibles
func getValidationMessage(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "This field is required"
	case "min":
		return "This field must be at least " + fe.Param() + " characters"
	case "max":
		return "This field must be at most " + fe.Param() + " characters"
	case "oneof":
		return "This field must be one of: " + fe.Param()
	case "datetime":
		return "This field must be a valid datetime in format: " + fe.Param()
	default:
		return "This field is invalid"
	}
}
