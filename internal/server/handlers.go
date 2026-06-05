package server

import (
	"net/http"

	"github.com/alizarazot/2026-i-dbm/internal/auth"
	"github.com/alizarazot/2026-i-dbm/internal/database"

	"github.com/labstack/echo/v5"
)

func handlerSignIn(jwtSecret []byte, authStore *database.AuthStore) echo.HandlerFunc {
	return func(c *echo.Context) error {
		userCredentials := struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}{}
		if err := c.Bind(&userCredentials); err != nil {
			return c.NoContent(http.StatusBadRequest)
		}

		user, err := authStore.VerifyCredentials(c.Request().Context(), userCredentials.Email, userCredentials.Password)
		if err != nil {
			return c.NoContent(http.StatusUnauthorized)
		}

		c.Logger().Info("!D:", "user", user)

		// TODO: Implement refresh token.
		token, err := auth.CreateJWTToken(jwtSecret, user)
		if err != nil {
			return c.NoContent(http.StatusInternalServerError)
		}

		return c.JSON(http.StatusOK, map[string]string{"token": token})
	}
}
