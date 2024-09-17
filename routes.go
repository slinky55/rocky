package rocky

import (
	"github.com/labstack/echo/v4"
	"github.com/slinky55/rocky/controllers"
)

func RegisterRoutes(e *echo.Echo) {
	account := e.Group("/account")
	account.POST("/register", controllers.AccountRegister)
	account.POST("/verify", controllers.AccountVerify)

	auth := e.Group("/auth")
	auth.POST("/challenge", controllers.AuthChallenge)
	auth.POST("/login", controllers.AuthLogin)

	e.GET("/signal", controllers.Signal)
}
