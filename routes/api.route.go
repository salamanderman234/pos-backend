package routes

import (
	"github.com/labstack/echo/v4"
	"github.com/salamanderman234/pos-backend/controllers"
)

func routeRegisterAPI(server *echo.Group) {
	// AUTH ROUTES
	authGroup := server.Group("/auth")
	authGroup.POST("/login", controllers.AuthLogin)
	authGroup.POST("/send_2f", controllers.AuthResendTwoFactor)
	authGroup.POST("/verify_2f", controllers.AuthVerifyTwoFactor)
	authGroup.POST("/verify", controllers.AuthVerifyUser)
	authGroup.POST("/verify/send", controllers.AuthVerifyUser)
	authGroup.POST("/reset", controllers.AuthResetPassword)
	authGroup.POST("/reset/send", controllers.AuthSendResetPassword)
	// END OF AUTH ROUTES
}
