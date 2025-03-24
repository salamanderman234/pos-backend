package services

import (
	"context"
	"time"

	"github.com/pquerna/otp/totp"
	"github.com/salamanderman234/pos-backend/config"
	"github.com/salamanderman234/pos-backend/helpers"
	"github.com/salamanderman234/pos-backend/models"
	"github.com/salamanderman234/pos-backend/repositories"
	"golang.org/x/crypto/bcrypt"
)

func UserFindUser(ctx context.Context, id any, preloads []string, selects ...string) (models.User, error) {
	user, err := repositories.UserFindByID(ctx, id, selects, preloads)
	return user, err
}

func UserGenerateKey(ctx context.Context, id any, purpose config.UserKeyPurposeEnum, until time.Duration) (string, error) {
	key := helpers.GenerateRandomString(6, helpers.NUMBER_CHARSET, helpers.UPPERCASE_CHARSET)
	selects := []string{
		"key",
		"key_valid_until",
		"key_purpose",
	}
	data := models.User{
		Key:           key,
		KeyValidUntil: time.Now().Add(until).Unix(),
		KeyPurpose:    string(purpose),
	}
	_, err := repositories.UserUpdate(ctx, id, data, selects)
	return key, err
}

func UserUpdate(ctx context.Context, id any, data models.User) (models.User, error) {
	selects := []string{
		"fullname",
		"full_address",
		"province",
		"city",
	}
	return repositories.UserUpdate(ctx, id, data, selects)
}

func UserEnableTwoFactor(ctx context.Context, id any, method config.TwoFactorMethodEnum) (string, error) {
	user, err := UserFindUser(ctx, id, []string{}, "id", "username", "secret")
	if err != nil {
		return "", err
	}
	selects := []string{
		"is_two_factor_enabled",
		"two_factor_method",
	}
	data := models.User{
		IsTwoFactorEnabled: true,
		TwoFactorMethod:    string(method),
	}
	secret := ""
	if method == config.TwoFactorEnum_GA && secret == "" {
		selects = append(selects, "secret")
		result, err := totp.Generate(totp.GenerateOpts{
			Issuer:      config.ApplicationName(),
			AccountName: user.Username,
		})
		if err != nil {
			return secret, err
		}
		secret = result.Secret()
		data.Secret = secret
	}
	_, err = repositories.UserUpdate(ctx, id, data, selects)
	return secret, err
}

func UserDisableTwoFactor(ctx context.Context, id any, password string) error {
	user, err := UserFindUser(ctx, id, []string{}, "id", "password")
	if err != nil {
		return err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return config.ErrInvalidCredentials
	}
	selects := []string{"is_two_factor_enabled"}
	data := models.User{IsTwoFactorEnabled: false}
	_, err = repositories.UserUpdate(ctx, id, data, selects)
	return err
}

func UserRemoveSecret(ctx context.Context, id any, password string) error {
	user, err := UserFindUser(ctx, id, []string{}, "id", "password")
	if err != nil {
		return err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return config.ErrInvalidCredentials
	}
	selects := []string{"is_two_factor_enabled", "two_factor_method", "secret"}
	data := models.User{IsTwoFactorEnabled: false, TwoFactorMethod: "", Secret: ""}
	_, err = repositories.UserUpdate(ctx, id, data, selects)
	return err
}
