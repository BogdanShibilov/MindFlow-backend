package jwt

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/bogdanshibilov/mindflowbackend/internal/mindflow/entity"
)

var (
	ErrUnknownClaims = errors.New("unknows claims")
)

type UserClaims struct {
	UserUuid string
	Email    string
	Roles    []string
	jwt.RegisteredClaims
}

func NewAccessToken(user *entity.User, secret string, duration time.Duration) (string, error) {
	claims := UserClaims{
		UserUuid: user.Uuid.String(),
		Email:    user.Email,
		Roles:    user.Roles,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
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
