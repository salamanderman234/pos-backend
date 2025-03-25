package middlewares

import (
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/salamanderman234/pos-backend/config"
	"github.com/salamanderman234/pos-backend/helpers"
	"github.com/salamanderman234/pos-backend/models"
	"github.com/salamanderman234/pos-backend/services"
)

func AuthVerifyMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		bearerToken := c.Request().Header.Get("Authorization")
		if bearerToken == "" {
			return helpers.HandleError(c, config.ErrInvalidToken)
		}
		token := strings.TrimPrefix(bearerToken, "Bearer ")
		claims, err := helpers.JWTParseToken(token, string(config.ApplicationKey()))
		if err != nil {
			return helpers.HandleError(c, err)
		}
		id := claims[config.AUTH_TOKEN_ID_KEY].(string)
		selects := []string{}
		preloads := []string{"Devices"}
		user, err := services.UserFindUser(ctx, id, preloads, selects...)
		if err != nil {
			return helpers.HandleError(c, config.ErrInvalidToken)
		}
		if err := services.AuthCheckUserSuspendBanState(user); err != nil {
			return helpers.HandleError(c, err)
		}
		c.Set(config.SESSION_USER_KEY, user)
		c.Set(config.SESSION_TOKEN_KEY, token)
		return next(c)
	}
}

func AuthVerifyUserDeviceMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		device := c.Get(config.SESSION_DEVICE_KEY).(string)
		user := c.Get(config.SESSION_USER_KEY).(models.User)
		if err := services.AuthCheckBannedDevice(ctx, user, device); err != nil {
			return helpers.HandleError(c, err)
		}
		return next(c)
	}
}

func AuthOnlyAdminMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		_, ok := c.Get(config.SESSION_USER_KEY).(models.User)
		if ok {
			return helpers.HandleError(c, config.ErrFailedPolicy)
		}
		return next(c)
	}
}

func AuthVerifiedOnlyMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		user := c.Get(config.SESSION_USER_KEY).(models.User)
		if user.VerifiedAt == 0 {
			return helpers.HandleError(c, config.ErrUpgradeRequired)
		}
		return next(c)
	}
}
