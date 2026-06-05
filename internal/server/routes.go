package server

import (
	"net/http"

	"github.com/labstack/echo/v5"
)

func addRoutes(e *echo.Echo) {
	e.GET("/", func(c *echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
}
