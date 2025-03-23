package config 

import (
	"http"
)

type Response struct {
	Status	int
	Message	string
	Data 	any
}

func (h Response) Error() string {
	return h.Message
}

var (
	ErrInvalidCredentials = Response{
		Status:		http.StatusUnauthorized,
		Message:	"Invalid credentials",
	},
	ErrUserBanned = Response {
		Status:		418,
		Message:	"User is currently banned",
	},
	ErrUserSuspended = Response {
		Status:		419,
		Message:	"User is currently suspended",
	},
	ErrInvalidKey = Response {
		Status:		http.StatusBadRequest,
		Message:	"Invalid key for this action",
	},
	ErrKeyExpired = Response {
		Status:		http.StatusBadRequest,
		Message:	"This key already expired, please issue a new one !",
	},
	ErrBadRequest = Response {
		Status:		http.StatusBadRequest,
		Message:	"Bad request",
	}
)
