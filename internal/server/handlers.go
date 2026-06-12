package server

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/alizarazot/2026-i-dbm/internal/auth"
	"github.com/alizarazot/2026-i-dbm/internal/database"
	"github.com/alizarazot/2026-i-dbm/internal/model"

	"github.com/labstack/echo/v5"
)

const authJWTCookieName = "auth-jwt"

func handlerSignIn(jwtSecret []byte, userStore *database.UserStore) echo.HandlerFunc {
	// TODO: Redirect if there is a valid auth cookie.
	return func(c *echo.Context) error {
		userCredentials := struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}{}
		if err := c.Bind(&userCredentials); err != nil {
			return c.NoContent(http.StatusBadRequest)
		}

		user, err := userStore.VerifyCredentials(c.Request().Context(), userCredentials.Email, userCredentials.Password)
		if errors.Is(err, database.ErrUserInvalidCredentials) {
			return c.NoContent(http.StatusUnauthorized)
		}
		if err != nil {
			return c.NoContent(http.StatusInternalServerError)
		}

		// TODO: Implement refresh token.
		token, expiration, err := auth.CreateJWTToken(jwtSecret, user)
		if err != nil {
			return c.NoContent(http.StatusInternalServerError)
		}

		c.SetCookie(&http.Cookie{
			Name:     authJWTCookieName,
			Value:    token,
			Path:     "/api",
			Expires:  expiration,
			Secure:   true,
			HttpOnly: true,
			SameSite: http.SameSiteLaxMode,
		})

		return c.NoContent(http.StatusOK)
	}
}

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
		var userRaw struct {
			Email           string `json:"email"`
			InitialPassword string `json:"initialPassword"`

			Role model.UserRole `json:"role"`
			Info model.UserInfo `json:"info"`
		}
		if err := echo.BindBody(c, &userRaw); err != nil {
			return echo.ErrBadRequest.Wrap(err)
		}

		user := model.User{
			Email: userRaw.Email,
			Role:  userRaw.Role,
			Info:  userRaw.Info,
		}
		if err := userStore.AddUser(c.Request().Context(), &user, userRaw.InitialPassword); err != nil {
			if errors.Is(err, database.ErrUserAlreadyExists) {
				return echo.NewHTTPError(http.StatusConflict, fmt.Sprintf("user %q already exists", user.Email)).Wrap(err)
			}

			return echo.ErrInternalServerError.Wrap(err)
		}

		return c.NoContent(http.StatusOK)
	}
}
