package server

import (
	"net/http"

	"github.com/alizarazot/2026-i-dbm/internal/database"

	"github.com/labstack/echo/v5"
)

func addFrontendRoutes(e *echo.Group) {
	e.GET("/", func(c *echo.Context) error {
		return c.HTML(http.StatusOK, "<body><a href='/sign-in'>Sign In</a></body>")
	})
}

func addPublicAPIRoutes(e *echo.Group, jwtSecret []byte, authStore *database.AuthStore) {
	e.POST("/auth", handlerSignIn(jwtSecret, authStore))
}

func addAPIRoutes(e *echo.Group) {
	e.GET("/ping", func(c *echo.Context) error {
		return c.String(http.StatusOK, "Pong!")
	})
}
