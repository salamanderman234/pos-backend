package helpers

import (
	"net/http"
	"slices"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/salamanderman234/pos-backend/config"
)

type RequestBSVConfig struct {
	Sanitize             bool
	Validate             bool
	SanitizeExceptFields []string
}

var (
	ONLY_SANITIZE_CONFIG     = RequestBSVConfig{Sanitize: true}
	ONLY_VALIDATE_CONFIG     = RequestBSVConfig{Validate: true}
	VALIDATE_SANITIZE_CONFIG = RequestBSVConfig{Sanitize: true, Validate: true}
)

func RequestBSV(c echo.Context, target any, cfg RequestBSVConfig) error {
	if err := c.Bind(target); err != nil {
		return config.ErrBadRequest
	}
	data := map[string]any{}
	TranslateStruct(target, &data)
	if cfg.Sanitize {
		data = RequestSanitizeForm(data, cfg.SanitizeExceptFields)
	}

	TranslateStruct(data, target)
	if cfg.Validate {
		if err := RequestValidateForm(target); err != nil {
			return err
		}
	}
	return nil
}

func RequestSanitizeForm(data map[string]any, skipSanitize []string) map[string]any {
	policy := config.Sanitizer()
	for key, value := range data {
		if slices.Contains(skipSanitize, key) {
			data[key] = value
		}
		valStr, ok := value.(string)
		if ok {
			clear := policy.Sanitize(valStr)
			data[key] = clear
		}
	}
	return data
}

func RequestValidateForm(target any) error {
	errs := config.Validator().Struct(target)
	if errs != nil {
		return errs
	}
	return nil
}

func RequestVerifyLimitCookie(c echo.Context, cookieName string) error {
	cookie, _ := c.Cookie(cookieName)
	if cookie != nil {
		return config.ErrTooManyRequest
	}
	return nil
}

func RequestGenerateLimitCookie(c echo.Context, cookieName string, limit time.Duration) {
	cookie := new(http.Cookie)
	cookie.Name = cookieName
	cookie.Value = GenerateRandomString(16)
	cookie.Expires = time.Now().Add(limit)
	cookie.HttpOnly = true
	cookie.Path = "/"

	c.SetCookie(cookie)
}
