package config

import (
	"net/http"
)

type Response struct {
	Status  int
	Message string
	Data    any
	Debug   any `json:"debug,omitempty"`
}

type FieldError struct {
	Field string `json:"field"`
	Error string `json:"error"`
}

func (h Response) Error() string {
	return h.Message
}

// registered error
var (
	ErrInvalidCredentials = Response{
		Status:  http.StatusUnauthorized,
		Message: "Invalid credentials",
	}
	ErrUserBanned = Response{
		Status:  http.StatusTeapot,
		Message: "User is currently banned",
	}
	ErrUserSuspended = Response{
		Status:  http.StatusLocked,
		Message: "User is currently suspended",
	}
	ErrInvalidKey = Response{
		Status:  http.StatusBadRequest,
		Message: "Invalid key for this action",
	}
	ErrExpiredKey = Response{
		Status:  http.StatusBadRequest,
		Message: "This key already expired, please issue a new one !",
	}
	ErrBadRequest = Response{
		Status:  http.StatusBadRequest,
		Message: "Bad request",
	}
	ErrNotFound = Response{
		Status:  http.StatusNotFound,
		Message: "Not found",
	}
	ErrConflict = Response{
		Status:  http.StatusConflict,
		Message: "Resource conflicted",
	}
	ErrInvalidToken = Response{
		Status:  http.StatusUnauthorized,
		Message: "Invalid access token",
	}
	ErrTooLarge = Response{
		Status:  http.StatusRequestEntityTooLarge,
		Message: "Too large",
	}
	ErrInvalidMimeType = Response{
		Status:  http.StatusUnsupportedMediaType,
		Message: "Invalid mime type for file",
	}
	ErrMethodNotAllowed = Response{
		Status:  http.StatusMethodNotAllowed,
		Message: "Method not allowed",
	}
	ErrDeviceBanned = Response{
		Status:  http.StatusLocked,
		Message: "Your device has been banned",
	}
	ErrFailedPolicy = Response{
		Status:  http.StatusForbidden,
		Message: "You are not allowed to access this resource",
	}
	ErrUpgradeRequired = Response{
		Status:  http.StatusUpgradeRequired,
		Message: "You'll need to upgrade your account status first",
	}
	ErrTooManyRequest = Response{
		Status:  http.StatusTooManyRequests,
		Message: "Please wait to make another request",
	}
)
