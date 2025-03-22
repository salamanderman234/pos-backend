package config 

import (
	"http"
)

type HttpError struct {
	Status	int
	Message	string
}

func (h HttpError) Error() string {
	return h.Message
}

var (
	ErrInvalidCredentials = HttpError{
		Status:		http.StatusUnauthorized,
		Message:	"Invalid credentials",
	},
	ErrUserBanned = HttpError {
		Status:		418,
		Message:	"User is currently banned",
	},
	ErrUserSuspended = HttpError {
		Status:		419,
		Message:	"User is currently suspended"
	}
)
