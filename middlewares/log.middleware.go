package middlewares

import (
	"github.com/labstack/echo/v4"
	"github.com/salamanderman234/pos-backend/config"
	"github.com/salamanderman234/pos-backend/models"
	"github.com/salamanderman234/pos-backend/services"
)

func LogRequestMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		next(c)
		status := c.Response().Status
		url := c.Request().URL.Path
		method := c.Request().Method
		userID := uint(0)
		userIntf := c.Get(config.SESSION_USER_KEY)
		if userIntf != nil {
			user := userIntf.(models.User)
			userID = user.ID
		}
		device := c.Get(config.SESSION_DEVICE_KEY).(string)
		ip := c.RealIP()
		go services.LogDispatchRequest(userID, device, ip, method, url, status)
		return nil
	}
}
