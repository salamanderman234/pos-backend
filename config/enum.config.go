package config

type UserKeyPurposeEnum string

const (
	UserKeyPurposeEnum_RESET_PASSWORD = "reset"
	UserKeyPurposeEnum_VERIFY         = "verify"
)

type TwoFactorMethodEnum string

const (
	TwoFactorEnum_EMAIL = "email"
	TwoFactorEnum_GA    = "authenticator"
)

type LogDriverEnum string

const (
	LogDriverEnum_DATABASE          = "database"
	LogDriverEnum_SERVICE           = "service"
	LogDriverEnum_EXTERNAL_DATABASE = "external_database"
)

type LogTypeEnum string

const (
	LogTypeEnum_FAILURE       = "FAIL"
	LogTypeEnum_LOGIN_ATTEMPT = "LOGIN ATTEMPT"
	LogTypeEnum_USER_ACTIVITY = "USER ACTIVITY"
	LogTypeEnum_REQUEST       = "REQUEST"
	LogTypeEnum_UPDATE_LEVEL  = "UPDATE LEVEL"
)
