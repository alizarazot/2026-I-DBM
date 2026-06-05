package server

import (
	"log/slog"

	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
)

func NewServer(logger *slog.Logger) *echo.Echo {
	e := echo.NewWithConfig(echo.Config{Logger: logger})

	addRoutes(e)

	e.Use(middleware.RequestLogger())

	return e
}
