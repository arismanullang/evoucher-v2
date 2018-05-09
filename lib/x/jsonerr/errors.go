package jsonerr

import (
	"net/http"
)

var (
	ErrFatal = &ErrorResponse{
		Status:  http.StatusInternalServerError,
		Code:    "ERR_FATAL",
		Message: "An error occurred when processing your request.",
	}
	ErrUnauthorized = &ErrorResponse{
		Status:  http.StatusUnauthorized,
		Code:    "ERR_UNAUTHORIZED",
		Message: "You're not authorized.",
	}
	ErrForbidden = &ErrorResponse{
		Status:  http.StatusForbidden,
		Code:    "ERR_FORBIDDEN",
		Message: "You do not have the required permission to access.",
	}
)
