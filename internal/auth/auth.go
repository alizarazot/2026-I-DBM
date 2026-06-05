package auth

import (
	"time"

	"github.com/alizarazot/2026-i-dbm/internal/model"

	"github.com/golang-jwt/jwt/v5"
)

type jwtCustomClaims struct {
	model.User
	jwt.RegisteredClaims
}

func CreateJWTToken(jwtSecret []byte, user *model.User) (string, error) {
	claims := &jwtCustomClaims{
		*user,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 2)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	t, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", err
	}

	return t, nil
}
