package request

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
}

// ValidateStruct validates any struct based on struct tags
func ValidateStruct(s any) error {
	err := validate.Struct(s)
	if err != nil {
		validationErrors := err.(validator.ValidationErrors)
		return fmt.Errorf("%s", formatValidationErrors(validationErrors))
	}
	return nil
}

func formatValidationErrors(errors validator.ValidationErrors) string {
	var message string
	for _, e := range errors {
		message += fmt.Sprintf("Field '%s' failed validation: %s; ", e.Field(), e.Tag())
	}
	return message
}
