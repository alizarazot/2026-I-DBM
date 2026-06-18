package server

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/alizarazot/2026-i-dbm/internal/database"
	"github.com/alizarazot/2026-i-dbm/internal/model"
	"github.com/labstack/echo/v5"
)

func handlerManagerListUsers(userStore *database.UserStore) echo.HandlerFunc {
	return func(c *echo.Context) error {
		const pageParam = "page"
		page, err := echo.QueryParam[uint](c, pageParam)
		if err != nil {
			return echo.ErrBadRequest.Wrap(fmt.Errorf("error reading param %q: %w", pageParam, err))
		}

		const limitParam = "limit"
		limit, err := echo.QueryParam[uint](c, limitParam)
		if err != nil {
			return echo.ErrBadRequest.Wrap(fmt.Errorf("error reading param %q: %w", limitParam, err))
		}
		if limit < 1 {
			return echo.ErrBadRequest.Wrap(fmt.Errorf("the provided limit (%d) must be >= 1", limit))
		}

		const roleParam = "role"
		roleRaw := c.QueryParam(roleParam)

		role := model.NewUserRole(roleRaw)
		if roleRaw != "" && role == model.UserRoleInvalid {
			return echo.ErrBadRequest.Wrap(fmt.Errorf("unknown role %q", roleRaw))
		}

		if role == model.UserRoleInvalid {
			users, totalUsers, err := userStore.GetUsers(c.Request().Context(), page, limit)
			if err != nil {
				return echo.ErrInternalServerError.Wrap(err)
			}

			return c.JSON(http.StatusOK, map[string]any{"users": users, "totalUsers": totalUsers})
		}

		users, totalUsers, err := userStore.GetUsersByRole(c.Request().Context(), role, page, limit)
		if err != nil {
			return echo.ErrInternalServerError.Wrap(err)
		}

		return c.JSON(http.StatusOK, map[string]any{"users": users, "totalUsers": totalUsers})
	}
}

func handlerManagerAddUser(userStore *database.UserStore) echo.HandlerFunc {
	return func(c *echo.Context) error {
		var data struct {
			User struct {
				Email string         `json:"email"`
				Role  model.UserRole `json:"role"`
				Info  model.UserInfo `json:"info"`
			} `json:"user"`
			InitialPassword string `json:"initialPassword"`
		}
		if err := echo.BindBody(c, &data); err != nil {
			return echo.ErrBadRequest.Wrap(err)
		}

		user := model.User{
			Email: data.User.Email,
			Role:  data.User.Role,
			Info:  data.User.Info,
		}
		if err := userStore.AddUser(c.Request().Context(), &user, data.InitialPassword); err != nil {
			if errors.Is(err, database.ErrUserAlreadyExists) {
				return echo.NewHTTPError(http.StatusConflict, fmt.Sprintf("user %q already exists", user.Email)).Wrap(err)
			}

			return echo.ErrInternalServerError.Wrap(err)
		}

		return c.NoContent(http.StatusOK)
	}
}

func handlerManagerEditUser(userStore *database.UserStore) echo.HandlerFunc {
	return func(c *echo.Context) error {
		var data struct {
			User struct {
				Email string         `json:"email"`
				Role  model.UserRole `json:"role"`
				Info  model.UserInfo `json:"info"`
			} `json:"user"`
		}
		if err := echo.BindBody(c, &data); err != nil {
			return echo.ErrBadRequest.Wrap(err)
		}

		user := model.User{
			Email: data.User.Email,
			Role:  data.User.Role,
			Info:  data.User.Info,
		}
		if err := userStore.EditUser(c.Request().Context(), &user); err != nil {
			return echo.ErrInternalServerError.Wrap(err)
		}

		return c.NoContent(http.StatusOK)
	}
}

func handlerManagerDeleteUser(userStore *database.UserStore) echo.HandlerFunc {
	return func(c *echo.Context) error {
		var data struct {
			Email string `query:"email"`
		}
		if err := c.Bind(&data); err != nil {
			return echo.ErrBadRequest.Wrap(err)
		}

		if err := userStore.DeleteUser(c.Request().Context(), data.Email); err != nil {
			if errors.Is(err, database.ErrUserNotExists) {
				return echo.ErrNotFound.Wrap(err)
			}

			return echo.ErrInternalServerError.Wrap(err)
		}

		return nil
	}
}
