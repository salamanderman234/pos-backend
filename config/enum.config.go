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
