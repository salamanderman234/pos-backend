package config

import "time"

const (
	// app
	APP_KEY_SIZE                = 16
	APP_WORKER_NUM              = 16
	APP_WORKER_POOL_BUFFER_SIZE = 50
	// auth
	AUTH_TWO_FACTOR_SEP = "|"
	AUTH_TOKEN_ID_KEY   = "param"
	AUTH_TOKEN_NAME_KEY = "name"
	// session
	SESSION_DEVICE_KEY = "device"
	SESSION_USER_KEY   = "user"
	SESSION_IP_KEY     = "ip"
	SESSION_TOKEN_KEY  = "token"
	// header
	HEADER_CSRF = "X-CSRF-TOKEN"
	// job
	JOB_SEND_MAIL_RETRY = 3
	JOB_LOG_RETRY       = 10
	// time
	TIME_JWT_EXPIRE = 3 * 24 * time.Hour
	TIME_TWO_FACTOR = 1 * time.Minute
	TIME_VERIFY_KEY = 24 * time.Hour
	TIME_RESET_KEY  = 30 * time.Minute
	TIME_LIMIT_SEND = 1 * time.Minute
	// cookie
	COOKIE_VERIFY_LIMIT_COOKIE = "session_verify"
	COOKIE_RESET_LIMIT_COOKIE  = "reset_header"
)
