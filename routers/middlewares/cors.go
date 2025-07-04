package middlewares

import "github.com/labstack/echo"

func ExampleMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		return next(c)
	}
}
