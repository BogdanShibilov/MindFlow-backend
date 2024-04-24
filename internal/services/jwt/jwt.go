package jwtservice

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type UserClaims struct {
	Uuid  string   `json:"uuid"`
	Email string   `json:"email"`
	Roles []string `json:"roles"`
	jwt.RegisteredClaims
}

func NewAccessToken(
	uuid string,
	email string,
	roles []string,
	secret string,
	duration time.Duration,
) (string, error) {
	issuedAt := time.Now()
	expiresAt := time.Now().Add(duration)
	claims := UserClaims{
		Uuid:  uuid,
		Email: email,
		Roles: roles,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(issuedAt),
			ExpiresAt: jwt.NewNumericDate(expiresAt),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func ParseJwtToken(tokenString string, secret []byte) (*UserClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &UserClaims{}, func(t *jwt.Token) (interface{}, error) {
		return secret, nil
	})

	switch {
	case token.Valid:
		if claims, ok := token.Claims.(*UserClaims); ok {
			return claims, nil
		}
		return nil, ErrUnknownClaims
	default:
		return nil, err
	}
}
