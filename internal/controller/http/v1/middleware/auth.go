package middleware

import (
	"errors"
	"log/slog"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	userrepo "github.com/bogdanshibilov/mindflowbackend/internal/repository/user"
	jwtservice "github.com/bogdanshibilov/mindflowbackend/internal/services/jwt"
	userservice "github.com/bogdanshibilov/mindflowbackend/internal/services/user"
)

func RequireJwt(secret string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tokenString, err := getAuthorizationToken(ctx)
		if err != nil || tokenString == "null" || tokenString == "undefined" {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		claims, err := jwtservice.ParseJwtToken(tokenString, []byte(secret))
		if err != nil {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		ctx.Set("claims", claims)

		ctx.Next()
	}
}

// Parses claims and sets values in context for: "uuid", "email", "roles"
// Must always go after RequireJwt middleware
func ParseClaimsIntoContext() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		claimsMaybe, exists := ctx.Get("claims")
		if !exists {
			ctx.AbortWithStatus(http.StatusForbidden)
			return
		}

		claims, ok := claimsMaybe.(*jwtservice.UserClaims)
		if !ok {
			ctx.AbortWithStatus(http.StatusForbidden)
			return
		}

		ctx.Set("uuid", claims.Uuid)
		ctx.Set("email", claims.Email)
		ctx.Set("roles", claims.Roles)

		ctx.Next()
	}
}

// Gets token from context and checks if its claims have admin permission
// Must always go after RequireJwt middleware
func RequireAdminPermission(userservice *userservice.Service, log *slog.Logger) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		const op = "middleware.auth.RequireAdminPermission"

		claimsMaybe, exists := ctx.Get("claims")
		if !exists {
			ctx.AbortWithStatus(http.StatusForbidden)
			return
		}

		claims, ok := claimsMaybe.(*jwtservice.UserClaims)
		if !ok {
			ctx.AbortWithStatus(http.StatusForbidden)
			return
		}

		isAdmin, err := userservice.IsAdmin(ctx, claims.Uuid)
		if err != nil {
			if errors.Is(err, userrepo.ErrUserNotFound) {
				log.Warn("non admin user tried to enter admin route")
			} else {
				log.Error("failed to check if user is admin", op, err)
			}
			ctx.AbortWithStatus(http.StatusForbidden)
			return
		}

		if !isAdmin {
			ctx.AbortWithStatus(http.StatusForbidden)
			return
		}

		ctx.Next()
	}
}

func getAuthorizationToken(ctx *gin.Context) (string, error) {
	tokenHeader := ctx.Request.Header.Get("authorization")
	tokenFields := strings.Fields(tokenHeader)
	if len(tokenFields) != 2 || tokenFields[0] != "Bearer" {
		return "", errors.New("invalid or absent Authorization header")
	}

	return tokenFields[1], nil
}
