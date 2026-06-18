package server

import (
	"net/http"
	"time"

	"github.com/alizarazot/2026-i-dbm/internal/auth"
	"github.com/alizarazot/2026-i-dbm/internal/database"
	"github.com/alizarazot/2026-i-dbm/internal/model"
	"github.com/golang-jwt/jwt/v5"

	"github.com/labstack/echo/v5"
)

func addPublicAPIRoutes(e *echo.Group, jwtSecret []byte, userStore *database.UserStore) {
	e.POST("/auth", handlerSignIn(jwtSecret, userStore))
}

func addAPIRoutes(e *echo.Group, userStore *database.UserStore) {
	// TODO: Sign out route.
	e.GET("/auth", func(c *echo.Context) error {
		token, err := echo.ContextGet[*jwt.Token](c, "user")
		if err != nil {
			return echo.ErrInternalServerError.Wrap(err)
		}
		return c.JSON(http.StatusOK, map[string]any{"user": auth.ExtractUser(token)})
	})
	e.GET("/sign-out", func(c *echo.Context) error {
		c.SetCookie(&http.Cookie{
			Name:     authJWTCookieName,
			Path:     "/api",
			Expires:  time.Unix(0, 0), // Delete cookie.
			Secure:   true,
			HttpOnly: true,
			SameSite: http.SameSiteLaxMode,
		})

		return c.NoContent(http.StatusOK)
	})

	manager := e.Group("/manager")
	manager.Use(middlewareRequireRole(model.UserRoleManager))
	manager.GET("/users", handlerManagerListUsers(userStore))
	manager.POST("/user", handlerManagerAddUser(userStore))
	manager.PUT("/user", handlerManagerEditUser(userStore))
	manager.DELETE("/user", handlerManagerDeleteUser(userStore))

	teacher := e.Group("/teacher")
	teacher.Use(middlewareRequireRole(model.UserRoleTeacher))
	teacher.GET("/ping", func(c *echo.Context) error {
		return c.String(http.StatusOK, "pong")
	})

	student := e.Group("/student")
	student.Use(middlewareRequireRole(model.UserRoleStudent))
	student.GET("/ping", func(c *echo.Context) error {
		return c.String(http.StatusOK, "pong")
	})
}
