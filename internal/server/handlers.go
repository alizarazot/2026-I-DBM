package server

import (
	"errors"
	"net/http"

	"github.com/alizarazot/2026-i-dbm/internal/auth"
	"github.com/alizarazot/2026-i-dbm/internal/database"
	"github.com/alizarazot/2026-i-dbm/internal/model"
	"github.com/golang-jwt/jwt/v5"

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

func handlerCommonAddCFC(userStore *database.UserStore, cfcStore *database.CFCStore) echo.HandlerFunc {
	return func(c *echo.Context) error {
		var data struct {
			Subject  string            `json:"subject"`
			Category model.CFCCategory `json:"category"`
			Details  string            `json:"details"`
		}
		if err := echo.BindBody(c, &data); err != nil {
			return echo.ErrBadRequest.Wrap(err)
		}

		token, err := echo.ContextGet[*jwt.Token](c, "user")
		if err != nil {
			return echo.ErrInternalServerError.Wrap(err)
		}
		user := auth.ExtractUser(token)

		cfc := model.CFC{
			Subject:  data.Subject,
			Category: data.Category,
			Details:  data.Details,
		}

		userID, err := userStore.GetUserID(c.Request().Context(), user.Email)
		if err != nil {
			return echo.ErrBadRequest.Wrap(err)
		}

		if err := cfcStore.AddCFC(c.Request().Context(), &cfc, userID); err != nil {
			return echo.ErrInternalServerError.Wrap(err)
		}

		return nil
	}
}
