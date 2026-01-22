package validation

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/gadz82/go-api-boilerplate/internal/domain"
)

// PlaygroundValidator implements domain.Validator using go-playground/validator
// This follows the Single Responsibility Principle - only handles validation logic
type PlaygroundValidator struct {
	validate *validator.Validate
}

// NewValidator creates a new PlaygroundValidator instance
// This is the constructor that will be used by the DI container
func NewValidator() domain.Validator {
	v := validator.New()

	// Register custom validation for UUID v4
	v.RegisterValidation("uuid4", validateUUID4)

	return &PlaygroundValidator{
		validate: v,
	}
}

// validateUUID4 is a custom validation function for UUID v4
func validateUUID4(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	if value == "" {
		return true // Empty is valid, use 'required' tag if field is mandatory
	}
	_, err := uuid.Parse(value)
	return err == nil
}

// Validate validates a struct and returns validation errors if any
func (pv *PlaygroundValidator) Validate(obj interface{}) domain.ValidationErrors {
	err := pv.validate.Struct(obj)
	if err == nil {
		return nil
	}

	return pv.translateErrors(err)
}

// ValidateField validates a single field value against a tag
func (pv *PlaygroundValidator) ValidateField(field interface{}, tag string) domain.ValidationErrors {
	err := pv.validate.Var(field, tag)
	if err == nil {
		return nil
	}

	return pv.translateErrors(err)
}

// translateErrors converts validator.ValidationErrors to domain.ValidationErrors
func (pv *PlaygroundValidator) translateErrors(err error) domain.ValidationErrors {
	var validationErrors domain.ValidationErrors

	if errs, ok := err.(validator.ValidationErrors); ok {
		for _, e := range errs {
			validationErrors = append(validationErrors, domain.ValidationError{
				Field:   toSnakeCase(e.Field()),
				Message: formatErrorMessage(e),
			})
		}
	}

	return validationErrors
}

// formatErrorMessage creates a human-readable error message
func formatErrorMessage(e validator.FieldError) string {
	field := toSnakeCase(e.Field())

	switch e.Tag() {
	case "required":
		return fmt.Sprintf("%s is required", field)
	case "uuid4":
		return fmt.Sprintf("%s must be a valid UUID v4", field)
	case "min":
		return fmt.Sprintf("%s must be at least %s characters", field, e.Param())
	case "max":
		return fmt.Sprintf("%s must be at most %s characters", field, e.Param())
	case "email":
		return fmt.Sprintf("%s must be a valid email address", field)
	case "url":
		return fmt.Sprintf("%s must be a valid URL", field)
	default:
		return fmt.Sprintf("%s failed validation: %s", field, e.Tag())
	}
}

// toSnakeCase converts a CamelCase string to snake_case
func toSnakeCase(s string) string {
	var result strings.Builder
	for i, r := range s {
		if i > 0 && r >= 'A' && r <= 'Z' {
			result.WriteRune('_')
		}
		result.WriteRune(r)
	}
	return strings.ToLower(result.String())
}
