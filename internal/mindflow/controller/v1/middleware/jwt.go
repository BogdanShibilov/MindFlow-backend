package middleware

import (
	"errors"
	"net/http"
	"strings"

	"github.com/bogdanshibilov/mindflowbackend/internal/mindflow/lib/jwt"
	"github.com/gin-gonic/gin"
)

var (
	ErrNoJwtHeader = errors.New("no authorization header")
)

func RequireJwt(secret string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tokenString, err := getAuthorizationToken(ctx)
		if err != nil {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		claims, err := jwt.ParseJwtToken(tokenString, []byte(secret))
		if err != nil {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		ctx.Set("userId", claims.UserUuid)
		ctx.Set("email", claims.Email)
		ctx.Set("roles", claims.Roles)

		ctx.Next()
	}
}

func getAuthorizationToken(ctx *gin.Context) (string, error) {
	tokenHeader := ctx.Request.Header.Get("Authorization")
	tokenFields := strings.Fields(tokenHeader)
	if len(tokenFields) != 2 || tokenFields[0] != "Bearer" {
		return "", ErrNoJwtHeader
	}

	return tokenFields[1], nil
}
