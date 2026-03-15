package pkg

import "errors"

// Custom error types for the application
var (
	ErrUserNotFound    = errors.New("user not found")
	ErrUserExists      = errors.New("user already exists")
	ErrInvalidPassword = errors.New("invalid password")
	ErrTokenExpired    = errors.New("token has expired")
	ErrTokenInvalid    = errors.New("invalid token")
	ErrUnauthorized    = errors.New("unauthorized access")
	ErrForbidden       = errors.New("forbidden access")
	ErrMovieNotFound   = errors.New("movie not found")
	ErrGenreNotFound   = errors.New("genre not found")
	ErrShowtimeNotFound = errors.New("showtime not found")
	ErrSeatNotAvailable = errors.New("seat not available")
	ErrReservationNotFound = errors.New("reservation not found")
	ErrReservationPast   = errors.New("cannot cancel past reservation")
	ErrAlreadyReserved   = errors.New("seat already reserved")
	ErrInvalidDate       = errors.New("invalid date")
	ErrDatabase          = errors.New("database error")
	ErrValidation        = errors.New("validation error")
	ErrRecordNotFound    = errors.New("record not found")
	ErrDuplicateEmail    = errors.New("email already exists")
)

// ValidationError represents a validation error with field details
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// NewValidationError creates a new validation error
func NewValidationError(field, message string) *ValidationError {
	return &ValidationError{Field: field, Message: message}
}

// Error implements the error interface
func (v *ValidationError) Error() string {
	return v.Message
}

// ValidationErrorList represents a list of validation errors
type ValidationErrorList struct {
	Errors []*ValidationError `json:"errors"`
}

// NewValidationErrorList creates a new validation error list
func NewValidationErrorList() *ValidationErrorList {
	return &ValidationErrorList{Errors: make([]*ValidationError, 0)}
}

// Add adds a validation error to the list
func (v *ValidationErrorList) Add(field, message string) {
	v.Errors = append(v.Errors, NewValidationError(field, message))
}

// Error implements the error interface
func (v *ValidationErrorList) Error() string {
	if len(v.Errors) == 0 {
		return "validation error"
	}
	return v.Errors[0].Error()
}
