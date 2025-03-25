package routes

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/salamanderman234/pos-backend/config"
	"github.com/salamanderman234/pos-backend/helpers"
	"github.com/salamanderman234/pos-backend/middlewares"
)

// middlewares that apply to api routes
func routeAPIGetDefaultMiddleware() []echo.MiddlewareFunc {
	return []echo.MiddlewareFunc{
		middleware.LoggerWithConfig(middleware.LoggerConfig{
			Format: "method=${method}, uri=${uri}, status=${status}\n",
		}),
		middleware.RateLimiterWithConfig(middlewares.RateLimitconfig),
		middleware.CORSWithConfig(middleware.CORSConfig{
			AllowOrigins: []string{"*"},
		}),
		middleware.BodyLimit("10M"),
		middlewares.UserAgentWhitelistMiddleware,
		middlewares.IPWhitelistMiddleware,
		middlewares.SessionRetrieveDeviceMiddleware,
		middlewares.LogRequestMiddleware,
	}
}

// middlwares that apply to web routes
func routeWebGetDefaultMiddleware() []echo.MiddlewareFunc {
	return []echo.MiddlewareFunc{}
}

// change default error handling for echo
func routeHandleError(server *echo.Echo) {
	server.HTTPErrorHandler = func(err error, c echo.Context) {
		if c.Response().Committed {
			return
		}

		if he, ok := err.(*echo.HTTPError); ok && he.Code == http.StatusNotFound {
			_ = helpers.HandleError(c, config.ErrNotFound)

			return
		}

		if he, ok := err.(*echo.HTTPError); ok && he.Code == http.StatusMethodNotAllowed {
			_ = helpers.HandleError(c, config.ErrMethodNotAllowed)

			return
		}
		server.DefaultHTTPErrorHandler(err, c)
	}
}

func RouteSetup(server *echo.Echo) {
	// register api routes
	version := config.ApplicationVersion()
	prefix := fmt.Sprintf("/api/v%s", version)
	apiGroup := server.Group(prefix, routeAPIGetDefaultMiddleware()...)
	routeRegisterAPI(apiGroup)
	// register web routes
	prefix = ""
	webGroup := server.Group("", routeWebGetDefaultMiddleware()...)
	routeRegisterWeb(webGroup)
	// handle error
	routeHandleError(server)
}
