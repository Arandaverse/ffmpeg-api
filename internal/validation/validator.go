package validation

import (
	"encoding/json"
	"fmt"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
}

// ValidationErrors represents a map of field names to their validation error messages
type ValidationErrors struct {
	Errors map[string]string `json:"errors"`
}

// Error implements the error interface
func (ve *ValidationErrors) Error() string {
	jsonBytes, err := json.Marshal(ve.Errors)
	if err != nil {
		return "validation failed"
	}
	return string(jsonBytes)
}

// Validate validates a struct using validator tags
func Validate(s interface{}) error {
	if err := validate.Struct(s); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			fieldErrors := make(map[string]string)
			for _, e := range validationErrors {
				fieldErrors[e.Field()] = formatError(e)
			}
			return &ValidationErrors{Errors: fieldErrors}
		}
		return err
	}
	return nil
}

// formatError formats a validation error into a human-readable message
func formatError(e validator.FieldError) string {
	switch e.Tag() {
	case "required":
		return "should not be empty"
	case "min":
		return fmt.Sprintf("should be at least %s characters long", e.Param())
	case "max":
		return fmt.Sprintf("should not be longer than %s characters", e.Param())
	case "email":
		return "should be a valid email address"
	case "url":
		return "should be a valid URL"
	case "oneof":
		return fmt.Sprintf("should be one of: %s", e.Param())
	default:
		return fmt.Sprintf("failed validation: %s", e.Tag())
	}
}
