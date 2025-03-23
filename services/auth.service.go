package services

import (
	"context"
	"bcrypt"	
	"time"
	"github.com/pquerna/otp/totp"

	"github.com/salamanderman234/pos-backend/models"
	"github.com/salamanderman234/pos-backend/repositories"
	"github.com/salamanderman234/pos-backend/config"
	"github.com/salamanderman234/pos-backend/helpers"
)

func AuthDecodeTwoFactorString(encoded string, usingTime bool) (string, string, string, error) {
	mappedVal, err := helpers.DecodeString(encoded)
	if len(mappedVal) < 2 {
		return "", "", "", config.ErrInvalidKey
	}

	username := mappedVal[0]
	validKey := mappedVal[1]
	
	timeStr := ""
	if usingTime {
		if len(mappedVal) < 3 {
			return "", "", "", config.ErrInvalidKey	
		}
		timeStr = mappedVal[2]
	}	
	return username, validKey, timeStr, nil
} 

func AuthEncodeTwoFactorString(username string, until ...time.Time) (string, string, error) {
	validKey := helpers.GenerateRandomString(6, helpers.UPPERCASE_CHARSET, helpers.NUMBER_CHARSET)
	arrFormat := []string{
		username,
		validKey,
	}
	if len(until) >= 1 {
		timeStr := until[0].Format("2025-03-23 16:02:03")
		arrFormat = append(arrFormat, timeStr)
	}
	result := helpers.EncodeString(arrFormat...)
	return result, validKey, nil
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
		"banned_at",
		"ban_reason",
		"suspended_at",
		"suspended_until",
		"suspend_reason",
		"verified_at",
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

func AuthVerififyTwoFactor(ctx context.Context, encoded string, key string) (string, error) {
	username, validKey, exp, err := AuthDecodeTwoFactorString(encoded, true)
	if err != nil {
		return "", config.ErrInvalidKey
	}
	untilParsed := time.Parse("2025-03-23 15:03:03", exp)
	if time.Now().After(untilParsed) {
		return "", config.ErrExpiredKey
	}

	selects := []string{
		"username",
		"email",
		"fullname",
		"banned_at",
		"ban_reason",
		"suspended_at",
		"suspend_reason",
		"suspended_until",
		"verified_at",
		"two_factor_method",
		"is_two_factor_enabled",
		"secret",
	}
	preloads := []string{"Notifications"}
	user, err := repositories.UserFindByUsername(ctx, username, selects, preloads)
	if err != nil {
		return "", config.ErrInvalidKey
	}	

	if err := AuthCheckUserSuspendBanState(user); err != nil {
		return "", err
	}
	if !user.IsTwoFactorEnabled {
		return "", config.ErrInvalidKey
	}
	method := user.TwoFactorMethod
	switch(method) {
	case config.TwoFactorEnum_EMAIL:
		if validKey != key {
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

func AuthResendTwoFactor(ctx context.Context, encoded string) (models.User, string, string, error) {
	username, _, _, err := AuthDecodeTwoFactorString(encoded, true)
	if err != nil {
		return models.User{}, "", "", err
	}
	selects := []string{"username", "two_factor_method", "is_two_factor_enabled"}
	preloads := []string{}
	user, err := repositories.UserFindByUsername(ctx, username, selects, preloads)
	if err != nil {
		return user, "", "", config.ErrInvalidKey
	}
	if !user.IsTwoFactorEnabled {
		return user, "", "", config.ErrInvalidKey
	}
	encodedResult, validKey, err := AuthEncodeTwoFactorString(username, time.Now().Add(config.TIME_TWO_FACTOR)
	
	return user, encodedResult, validKey, err
}

func AuthVerifyUser(ctx context.Context, key string, username string) (models.User, error) {
	nowUnix := time.Now().Unix()
	selects := []string{
		"id", 
		"username", 
		"key", 
		"key_valid_until",
		"verified_at",
	}
	preloads := []string{}
	user, err := models.UserFindByUsername(ctx, username, selects, preloads)
	if err != nil {
		return user, config.ErrInvalidKey
	}
	if user.VerifiedAt != 0 {
		return user, config.ErrInvalidKey
	}
	if key != user.Key {
		return user, config.ErrInvalidKey
	}
	if nowUnix > user.KeyValidUntil {
		return user, config.ErrExpiredKey
	}
	
	id := user.ID
	selects = []string{
		"key",
		"key_valid_until",
		"verified_at",
	}
	data := models.User{
		VerifiedAt: 	nowUnix,
		Key:		"",
		KeyValidUntil:	0,
	}
	return repositories.UserUpdate(id, data, selects...)
}

func AuthResetPassword(ctx context.Context, username string, code string, newPassword string) error {
	nowUnix := time.Now().Unix()
	selects := []string{
		"username",
		"key",
		"key_valid_until",
		"key_purpose",
	}

	preloads := []string{"OldPasswords"}
	user, err := repositories.UserFindByUsername(ctx, username, selects, preloads)
	if err != nil {
		return err
	}	
	if user.KeyPurpose != config.UserKeyPurposeEnum_RESET_PASSWORD {
		return config.ErrInvalidKey
	}
	if code != user.Key {
		return config.ErrInvalidKey
	}
	if nowUnix > user.KeyValidUntil {
		return config.ErrExpiredKey
	}
	
}
