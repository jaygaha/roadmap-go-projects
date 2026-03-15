package pkg

import (
	"regexp"
	"strings"
)

// EmailValidator validates email format
type EmailValidator struct{}

// Validate checks if the email is valid
func (e *EmailValidator) Validate(email string) bool {
	// Basic email regex pattern
	pattern := `^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`
	matched, _ := regexp.MatchString(pattern, email)
	return matched
}

// NewEmailValidator creates a new email validator
func NewEmailValidator() *EmailValidator {
	return &EmailValidator{}
}

// ValidateEmail validates an email address
func ValidateEmail(email string) bool {
	return NewEmailValidator().Validate(email)
}

// IsRequired checks if a string is not empty
func IsRequired(value string) bool {
	return strings.TrimSpace(value) != ""
}

// MinLength checks if string meets minimum length
func MinLength(value string, min int) bool {
	return len(value) >= min
}

// MaxLength checks if string exceeds maximum length
func MaxLength(value string, max int) bool {
	return len(value) <= max
}

// IsNumeric checks if string contains only digits
func IsNumeric(value string) bool {
	_, err := regexp.MatchString(`^\d+$`, value)
	return err == nil
}
