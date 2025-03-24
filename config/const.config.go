package config

import "time"

const (
	// app
	APP_KEY_SIZE = 16
	// auth
	AUTH_TWO_FACTOR_SEP = "|"
	AUTH_TOKEN_ID_KEY   = "param"
	AUTH_TOKEN_NAME_KEY = "name"
	// time
	TIME_JWT_EXPIRE = 3 * 24 * time.Hour
	TIME_TWO_FACTOR = 1 * time.Minute
)
