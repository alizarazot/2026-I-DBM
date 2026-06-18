package server

import (
	"fmt"

	"github.com/alizarazot/2026-i-dbm/internal/auth"
	"github.com/alizarazot/2026-i-dbm/internal/model"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v5"
)

func middlewareRequireRole(role model.UserRole) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			token, err := echo.ContextGet[*jwt.Token](c, "user")
			if err != nil {
				return echo.ErrUnauthorized.Wrap(err)
			}

			user := auth.ExtractUser(token)

			if user.Role != role {
				return echo.ErrUnauthorized.Wrap(fmt.Errorf("this endpoint needs role %q but %q was provided", role, user.Role))
			}

			return next(c)
		}
	}
}
