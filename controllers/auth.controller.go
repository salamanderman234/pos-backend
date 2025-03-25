package controllers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/salamanderman234/pos-backend/config"
	"github.com/salamanderman234/pos-backend/forms"
	"github.com/salamanderman234/pos-backend/helpers"
	"github.com/salamanderman234/pos-backend/response"
	"github.com/salamanderman234/pos-backend/services"
)

func AuthLogin(c echo.Context) error {
	ctx := c.Request().Context()
	device := c.Get("device").(string)
	ip := c.RealIP()

	form := forms.FormLogin{}
	if err := helpers.RequestBSV(c, &form, helpers.VALIDATE_SANITIZE_CONFIG); err != nil {
		return helpers.HandleError(c, err)
	}
	user, isTwoFactor, token, err := services.AuthLogin(ctx, form.Username, form.Password)
	go services.LogDispatchLoginAttempt(user.ID, device, ip, err == nil, isTwoFactor, token)
	if err != nil {
		return helpers.HandleError(c, err)
	}
	if isTwoFactor {
		method := user.TwoFactorMethod
		until := time.Now().Add(config.TIME_TWO_FACTOR)
		encoded, validKey, err := services.AuthEncodeTwoFactorString(form.Username, until)
		if err != nil {
			return helpers.HandleError(c, err)
		}
		if method == string(config.TwoFactorEnum_EMAIL) {
			go services.MailSendTwoFactor(user.Email, validKey)
		}
		payload := map[string]any{
			"seed":   encoded,
			"method": method,
		}
		return c.JSON(http.StatusOK, payload)
	}
	userResp := response.UserResponse{}
	helpers.TranslateStruct(user, &userResp)
	payload := map[string]any{
		"token": token,
		"user":  userResp,
	}

	return c.JSON(http.StatusOK, config.Response{
		Status:  http.StatusOK,
		Message: "OK",
		Data:    payload,
	})
}

func AuthResendTwoFactor(c echo.Context) error {
	ctx := c.Request().Context()
	form := forms.FormResendTwoFactor{}
	if err := helpers.RequestBSV(c, &form, helpers.ONLY_VALIDATE_CONFIG); err != nil {
		return helpers.HandleError(c, err)
	}
	user, encoded, validKey, err := services.AuthResendTwoFactor(ctx, form.Seed)
	if err != nil {
		return helpers.HandleError(c, err)
	}
	go services.MailSendTwoFactor(user.Email, validKey)
	payload := map[string]any{
		"seed":   encoded,
		"method": user.TwoFactorMethod,
	}
	return c.JSON(http.StatusOK, config.Response{
		Status:  http.StatusOK,
		Message: "OK",
		Data:    payload,
	})
}

func AuthVerifyTwoFactor(c echo.Context) error {
	device := c.Get("device").(string)
	ip := c.RealIP()
	ctx := c.Request().Context()
	form := forms.FormVerifyTwoFactor{}
	exceptSanitize := []string{"seed"}
	if err := helpers.RequestBSV(c, &form, helpers.RequestBSVConfig{Sanitize: true, Validate: true, SanitizeExceptFields: exceptSanitize}); err != nil {
		return helpers.HandleError(c, err)
	}
	user, token, err := services.AuthVerififyTwoFactor(ctx, form.Seed, form.Code)
	go services.LogDispatchLoginAttempt(user.ID, device, ip, err == nil, true, token)
	if err != nil {
		return helpers.HandleError(c, err)
	}
	userResp := response.UserResponse{}
	helpers.TranslateStruct(user, &userResp)
	payload := map[string]any{
		"token": token,
		"user":  userResp,
	}

	return c.JSON(http.StatusOK, config.Response{
		Status:  http.StatusOK,
		Message: "OK",
		Data:    payload,
	})
}

func AuthVerifyUser(c echo.Context) error {
	ctx := c.Request().Context()
	form := forms.FormVerifyUser{}
	if err := helpers.RequestBSV(c, &form, helpers.VALIDATE_SANITIZE_CONFIG); err != nil {
		return helpers.HandleError(c, err)
	}
	user, err := services.AuthVerifyUser(ctx, form.Key, form.Username)
	if err != nil {
		return helpers.HandleError(c, err)
	}
	go services.LogDispatchUserActivity(
		user.ID,
		fmt.Sprintf("User %s is successfully verify their account", user.Username),
	)

	return c.JSON(http.StatusOK, config.Response{
		Status:  http.StatusOK,
		Message: "OK",
	})
}

func AuthResendVerifyEmail(c echo.Context) error {
	if err := helpers.RequestVerifyLimitCookie(c, config.COOKIE_VERIFY_LIMIT_COOKIE); err != nil {
		return helpers.HandleError(c, err)
	}
	ctx := c.Request().Context()
	form := forms.FormResendEmail{}
	if err := helpers.RequestBSV(c, &form, helpers.VALIDATE_SANITIZE_CONFIG); err != nil {
		return helpers.HandleError(c, err)
	}
	user, err := services.UserFindUserByUsername(ctx, form.Username, []string{}, []string{"id", "email", "username"}...)
	if err == nil {
		key, err := services.UserGenerateKey(ctx, user.ID, config.UserKeyPurposeEnum_VERIFY, config.TIME_VERIFY_KEY)
		if err == nil {
			go services.MailSendVerify(user.Email, user.Username, key)
			go services.LogDispatchUserActivity(
				user.ID,
				fmt.Sprintf("User %s requesting a new key to verify their account", user.Username),
			)
		}
	}
	helpers.RequestGenerateLimitCookie(c, config.COOKIE_VERIFY_LIMIT_COOKIE, config.TIME_LIMIT_SEND)
	return c.JSON(http.StatusOK, config.Response{
		Status:  http.StatusOK,
		Message: "OK",
	})
}

func AuthResetPassword(c echo.Context) error {
	ctx := c.Request().Context()
	form := forms.FormResetPassword{}
	if err := helpers.RequestBSV(c, &form, helpers.VALIDATE_SANITIZE_CONFIG); err != nil {
		return helpers.HandleError(c, err)
	}
	user, err := services.AuthResetPassword(ctx, form.Username, form.Key, form.NewPassword)
	if err != nil {
		return helpers.HandleError(c, err)
	}
	go services.LogDispatchUserActivity(
		user.ID,
		fmt.Sprintf("User %s is successfully change password for their", user.Username),
	)

	return c.JSON(http.StatusOK, config.Response{
		Status:  http.StatusOK,
		Message: "OK",
	})
}

func AuthSendResetPassword(c echo.Context) error {
	if err := helpers.RequestVerifyLimitCookie(c, config.COOKIE_RESET_LIMIT_COOKIE); err != nil {
		return helpers.HandleError(c, err)
	}
	ctx := c.Request().Context()
	form := forms.FormResendEmail{}
	if err := helpers.RequestBSV(c, &form, helpers.VALIDATE_SANITIZE_CONFIG); err != nil {
		return helpers.HandleError(c, err)
	}
	user, err := services.UserFindUserByUsername(ctx, form.Username, []string{}, []string{"id", "email", "username"}...)
	if err == nil {
		key, err := services.UserGenerateKey(ctx, user.ID, config.UserKeyPurposeEnum_RESET_PASSWORD, config.TIME_RESET_KEY)
		if err == nil {
			go services.MailSendResetPassword(user.Email, user.Username, key)
			go services.LogDispatchUserActivity(
				user.ID,
				fmt.Sprintf("User %s requesting a new key to reset their password", user.Username),
			)
		}
	}
	helpers.RequestGenerateLimitCookie(c, config.COOKIE_RESET_LIMIT_COOKIE, config.TIME_LIMIT_SEND)
	return c.JSON(http.StatusOK, config.Response{
		Status:  http.StatusOK,
		Message: "OK",
	})
}
