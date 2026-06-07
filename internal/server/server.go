package server

import (
	"log/slog"

	"github.com/alizarazot/2026-i-dbm/internal/database"

	echojwt "github.com/labstack/echo-jwt/v5"
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
)

func NewServer(logger *slog.Logger, jwtSecret []byte, authStore *database.AuthStore) *echo.Echo {
	e := echo.NewWithConfig(echo.Config{Logger: logger})
	e.Use(middleware.RequestLogger())

	frontend := e.Group("")
	addFrontendRoutes(frontend)

	apiPublic := e.Group("/api")
	addPublicAPIRoutes(apiPublic, jwtSecret, authStore)

	api := e.Group("/api")
	api.Use(echojwt.WithConfig(echojwt.Config{
		SigningKey:  jwtSecret,
		TokenLookup: "cookie:" + authJWTCookieName,
	}))
	addAPIRoutes(api)

	return e
}
