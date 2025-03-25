package middlewares

import (
	"github.com/labstack/echo/v4"
	"github.com/salamanderman234/pos-backend/config"
)

func SessionRetrieveDeviceMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		userAgent := c.Request().Header.Get("User-Agent")
		c.Set(config.SESSION_DEVICE_KEY, userAgent)
		return next(c) // Continue to the next handler
	}
}
