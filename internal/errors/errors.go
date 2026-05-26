package errors

import (
	"fmt"
	"net/http"
	"strconv"
)

type APIError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details any    `json:"details,omitempty"`
	Status  int    `json:"-"`
}

type RateLimitError struct {
	APIError
	RetryAfter int `json:"retry_after"`
}

func (e *APIError) Error() string {
	return e.Message
}

func NotFound(message string) *APIError {
	return &APIError{
		Code:    CodeNotFound,
		Message: message,
		Status:  http.StatusNotFound,
	}
}

func Unauthorized(message string) *APIError {
	return &APIError{
		Code:    CodeUnauthorized,
		Message: message,
		Status:  http.StatusUnauthorized,
	}
}

func Forbidden(message string) *APIError {
	return &APIError{
		Code:    CodeForbidden,
		Message: message,
		Status:  http.StatusForbidden,
	}
}

func Conflict(message string) *APIError {
	return &APIError{
		Code:    CodeConflict,
		Message: message,
		Status:  http.StatusConflict,
	}
}

func InternalServerError(err error) *APIError {
	return &APIError{
		Code:    CodeInternal,
		Message: "Internal server error",
		Details: err.Error(),
		Status:  http.StatusInternalServerError,
	}
}

func TooManyRequests(ra int) *RateLimitError {
	return &RateLimitError{
		APIError: APIError{
			Code:    CodeTooManyRequests,
			Message: "Rate limit exceeded",
			Details: fmt.Sprintf("Too many requests. Please try again in %s seconds.", strconv.Itoa(ra)),
			Status:  http.StatusTooManyRequests,
		},
		RetryAfter: ra,
	}
}

func ValidationError(details interface{}) *APIError {
	return &APIError{
		Code:    CodeValidation,
		Message: "Validation failed",
		Details: details,
		Status:  http.StatusBadRequest,
	}
}

func FromGinValidation(err error) *APIError {
	// if validationErrs, ok := err.(validator.ValidationErrors); ok {
	// 	return ValidationError(appvalidator.TranslateValidationErrors(validationErrs))
	// }

	return &APIError{
		Code: CodeValidation,
		Message: "Invalid request data format",
		Details: err.Error(),
		Status: http.StatusBadRequest,
	}
}