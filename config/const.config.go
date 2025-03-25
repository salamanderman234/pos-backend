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
	// job
	JOB_SEND_MAIL_RETRY = 3
	JOB_LOG_RETRY       = 10
	// time
	TIME_JWT_EXPIRE = 3 * 24 * time.Hour
	TIME_TWO_FACTOR = 1 * time.Minute
)
