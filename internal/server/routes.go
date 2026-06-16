package server

import (
	"net/http"

	"github.com/alizarazot/2026-i-dbm/internal/database"
	"github.com/alizarazot/2026-i-dbm/internal/model"

	"github.com/labstack/echo/v5"
)

func addPublicAPIRoutes(e *echo.Group, jwtSecret []byte, userStore *database.UserStore) {
	e.POST("/auth", handlerSignIn(jwtSecret, userStore))
}

func addAPIRoutes(e *echo.Group, userStore *database.UserStore) {
	// TODO: Sign out route.
	e.GET("/auth", func(c *echo.Context) error {
		return c.NoContent(http.StatusOK)
	})

	manager := e.Group("/manager")
	manager.Use(middlewareRequireRole(model.UserRoleManager))
	manager.GET("/users", handlerManagerListUsers(userStore))
	manager.POST("/user", handlerManagerAddUser(userStore))

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
