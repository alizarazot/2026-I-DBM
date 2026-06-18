package server

import (
	"log/slog"

	"github.com/alizarazot/2026-i-dbm/frontend"
	"github.com/alizarazot/2026-i-dbm/internal/auth"
	"github.com/alizarazot/2026-i-dbm/internal/database"
	"github.com/golang-jwt/jwt/v5"

	echojwt "github.com/labstack/echo-jwt/v5"
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
)

func NewServer(
	logger *slog.Logger,
	jwtSecret []byte,
	userStore *database.UserStore,
	cfcStore *database.CFCStore,
) *echo.Echo {
	e := echo.NewWithConfig(echo.Config{Logger: logger})
	e.Use(middleware.RequestLogger())

	// Frontend SPA support.
	f := e.Group("")
	f.Use(middleware.StaticWithConfig(middleware.StaticConfig{
		Root:       ".",
		Filesystem: frontend.Files,
		HTML5:      true,
	}))
	f.RouteNotFound("/*", func(c *echo.Context) error {
		return c.FileFS(frontend.SPAFallbackFile, frontend.Files)
	})

	apiPublic := e.Group("/api")
	addPublicAPIRoutes(apiPublic, jwtSecret, userStore)

	api := e.Group("/api")
	api.Use(echojwt.WithConfig(echojwt.Config{
		SigningKey:  jwtSecret,
		TokenLookup: "cookie:" + authJWTCookieName,
		NewClaimsFunc: func(c *echo.Context) jwt.Claims {
			return auth.JWTClaims()
		},
	}))
	addAPIRoutes(api, userStore, cfcStore)

	return e
}
