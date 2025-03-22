package services

import (
	"context"
	"bcrypt"	
	"time"
	"github.com/pquerna/otp/totp"

	"github.com/salamanderman234/pos-backend/models"
	"github.com/salamanderman234/pos-backend/repositories"
	"github.com/salamanderman234/pos-backend/config"
)

func AuthGenerateTwoFactorEncodedString(user models.User) (string, error) {
	return "", nil
}

func AuthGenerateToken(user models.User) (string, error) {
	token, err := helpers.JWTGenerateTokenFromUser(user, config.TIME_JWT_EXPIRE)
	if err != nil {
		return token, err
	}
	return token, nil
}

func AuthCheckUserSuspendBanState(user models.User) error {
	if user.BannedAt != 0 {
		return config.ErrUserBanned
	}
	if user.SuspendedAt != 0 {
		return config.ErrUserSuspended
	}
	return nil
}

func AuthLogin(ctx context, username string, password string) (models.User, bool, string, error) {
	selects := []string{
		"username",
		"email",
		"password",
		"fullname",
		"avatar",
		"banned_at",
		"ban_reason",
		"suspended_at",
		"suspended_until",
		"suspend_reason",
		"activated_at",
		"two_factor_method",
		"is_two_factor_enabled",
	}
	preloads := []string{"Notifications"}
	user, err := repositories.UserFindUsername(ctx, username, selects)
	if err != nil {
		return user, false, "", config.ErrInvalidCredentials
	}

	hashed := user.Password
	if err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(password)); err != nil {
		return user, false, "", config.ErrInvalidCredentials
	}

	if err := AuthCheckUserSuspendBanState(user); err != nil {
		return user, false, "", err
	}
	
	if user.IsTwoFactorEnabled {
		return user, true, "", nil
	}

	token, err := AuthGenerateToken(user)
	if err != nil {
		return user, false, "", err
	}
	return user, false, token, nil
}

func VerififyTwoFactor(ctx context.Context, encoded string, key string) (string, error) {
	mappedVal, err := helpers.DecodeString(encoded)
	if err != nil {
		return "", err
	}
	if len(mappedVal) < 2 {
		return "", config.ErrBadRequest
	}
	method := mappedVal[0]
	username := mappedVal[1]
	
	selects := []string{
	
	}
	preloads := []string{}
	user, err := repositories.UserFindByUsername(ctx, username, selects, preloads)
	if err != nil {
		return "", config.ErrInvalidKey
	}	

	if err := AuthCheckUserSuspendBanState(user); err != nil {
		return "", err
	}

	switch(method) {
	case config.TwoFactorEnum_EMAIL:
		if len(mappedVal) < 4 {
			return "", config.ErrBadRequest
		}
		validKey := mappedVal[2]
		until := mappedVal[3]

		parsedUntil, _ := time.Parse("2025-03-20 16:03:04", until)

		if validKey != key || time.Now().After(parsedUntil) {
			return "", config.ErrInvalidKey		
		}
	case config.TwoFactorEnum_GA:
		secret := user.Secret
		if !totp.Validate(key, secret) {
			return "", config.ErrInvalidKey
		}	
	default:
		return "", config.ErrInvalidKey
	}
	
	token, err := AuthGenerateToken(user)
	if err != nil {
		return "", config.ErrInvalidKey
	}

	return token, nil
} 

func ResendTwoFactor(ctx context.Context, encoded string) error {

}
