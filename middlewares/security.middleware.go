package middlewares

import (
	"net/http"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/salamanderman234/pos-backend/config"
	"github.com/salamanderman234/pos-backend/helpers"
	"golang.org/x/time/rate"
)

func UserAgentWhitelistMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		agent := c.Request().Header.Get("User-Agent")
		accepted := "curl/"
		if !strings.Contains(agent, accepted) {
			return helpers.HandleError(c, config.ErrNotFound)
		}
		return next(c)
	}
}

func IPWhitelistMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		ip := c.RealIP()
		// get allowed ip from config
		allowedIps := []string{}
		isAllowed := false
		for _, allowedIp := range allowedIps {
			if ip == allowedIp || allowedIp == "*" {
				isAllowed = true
				break
			}
		}
		if !isAllowed {
			return helpers.HandleError(c, config.ErrNotFound)
		}
		return next(c)
	}
}

// list url that ignore rate limit
var ignoreRateLimit = []string{}

var RateLimitconfig = middleware.RateLimiterConfig{
	Skipper: func(c echo.Context) bool {
		for _, route := range ignoreRateLimit {
			if c.Path() == route {
				return false
			}
		}
		return true
	},
	Store: middleware.NewRateLimiterMemoryStoreWithConfig(
		middleware.RateLimiterMemoryStoreConfig{Rate: rate.Limit(10), Burst: 50, ExpiresIn: 10 * time.Minute},
	),
	IdentifierExtractor: func(ctx echo.Context) (string, error) {
		id := ctx.RealIP()
		return id, nil
	},
	ErrorHandler: func(context echo.Context, err error) error {
		return context.JSON(http.StatusForbidden, config.Response{
			Status:  http.StatusForbidden,
			Message: "Access from your IP address is not allowed.",
		})
	},
	DenyHandler: func(context echo.Context, identifier string, err error) error {
		return context.JSON(http.StatusTooManyRequests, config.Response{
			Status:  http.StatusTooManyRequests,
			Message: "You have made too many requests. Please try again later.",
		})
	},
}
