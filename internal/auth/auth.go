package auth

import (
	"time"

	"github.com/alizarazot/2026-i-dbm/internal/model"

	"github.com/golang-jwt/jwt/v5"
)

const jwtExpiration = time.Hour * 24 * 2

type jwtCustomClaims struct {
	model.User
	jwt.RegisteredClaims
}

func CreateJWTToken(jwtSecret []byte, user *model.User) (string, time.Time, error) {
	expiration := time.Now().Add(jwtExpiration)
	claims := &jwtCustomClaims{
		*user,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiration),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	t, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", time.Time{}, err
	}

	return t, expiration, nil
}

func ExtractUser(token *jwt.Token) *model.User {
	return &token.Claims.(*jwtCustomClaims).User
}

func JWTClaims() jwt.Claims {
	return &jwtCustomClaims{}
}
