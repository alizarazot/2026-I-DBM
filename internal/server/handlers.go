package server

import (
	"errors"
	"net/http"

	"github.com/alizarazot/2026-i-dbm/internal/auth"
	"github.com/alizarazot/2026-i-dbm/internal/database"

	"github.com/labstack/echo/v5"
)

const authJWTCookieName = "auth-jwt"

func handlerSignIn(jwtSecret []byte, authStore *database.AuthStore) echo.HandlerFunc {
	// TODO: Redirect if there is a valid auth cookie.
	return func(c *echo.Context) error {
		userCredentials := struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}{}
		if err := c.Bind(&userCredentials); err != nil {
			return c.NoContent(http.StatusBadRequest)
		}

		user, err := authStore.VerifyCredentials(c.Request().Context(), userCredentials.Email, userCredentials.Password)
		if errors.Is(err, database.ErrAuthInvalidCredentials) {
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
