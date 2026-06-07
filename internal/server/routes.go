package server

import (
	"net/http"

	"github.com/alizarazot/2026-i-dbm/internal/database"

	"github.com/labstack/echo/v5"
)

func addPublicAPIRoutes(e *echo.Group, jwtSecret []byte, authStore *database.AuthStore) {
	e.POST("/auth", handlerSignIn(jwtSecret, authStore))
}

func addAPIRoutes(e *echo.Group) {
	e.GET("/auth", func(c *echo.Context) error {
		return c.NoContent(http.StatusOK)
	})
}
