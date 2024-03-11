package jwt

import (
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/bogdanshibilov/mindflowbackend/internal/mindflow/entity"
)

func NewAccessToken(user *entity.User, secret string, duration time.Duration) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["uuid"] = user.Uuid
	claims["email"] = user.Email
	claims["expires_at"] = time.Now().Add(duration).Unix()

	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
