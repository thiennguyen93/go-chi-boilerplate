package routers

import (
	example "thiennguyen.dev/welab-healthcare-app/controllers"

	"github.com/labstack/echo/v4"
)

// RegisterRoutes add all routing list here automatically get main router
func RegisterRoutes(app *echo.Echo) {
	// Creating a New Router
	apiRouter := app.Group("/api/v1")
	apiRouter.GET("/test", example.GetData)
}
