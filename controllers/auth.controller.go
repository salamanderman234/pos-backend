package controllers

import (
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
	ip := c.Get("ip").(string)

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

	return c.JSON(http.StatusOK, payload)
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
	return c.JSON(http.StatusOK, payload)
}

func AuthVerifyTwoFactor(c echo.Context) error {
	return nil
}
