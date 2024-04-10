package routes

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/bogdanshibilov/mindflowbackend/internal/mindflow/controller/v1/dto"
	"github.com/bogdanshibilov/mindflowbackend/internal/mindflow/services/auth"
	"github.com/bogdanshibilov/mindflowbackend/internal/mindflow/storage"
)

type AuthRoutes struct {
	log  *slog.Logger
	auth *auth.Auth
}

func NewAuthRoutes(handler *gin.RouterGroup, log *slog.Logger, auth *auth.Auth) {
	r := &AuthRoutes{
		log:  log,
		auth: auth,
	}

	authHandler := handler.Group("/auth")
	{
		authHandler.POST("/register", r.Register)
		authHandler.POST("/login", r.Login)
	}
}

func (r *AuthRoutes) Register(ctx *gin.Context) {
	const op = "AuthRoutes.Register"

	var req *dto.RegisterRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		r.log.Error("failed to bind json data", op, err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON"})
		return
	}

	err := r.auth.RegisterNewUser(ctx, req.Email, req.Password)
	if err != nil {
		r.log.Error("failed to create a new user", op, err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create a new user"})
		return
	}

	ctx.Status(http.StatusCreated)
}

func (r *AuthRoutes) Login(ctx *gin.Context) {
	const op = "AuthRoutes.Login"

	var req *dto.LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		r.log.Error("failed to bind json data", op, err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON"})
		return
	}

	token, err := r.auth.Login(ctx, req.Email, req.Password)
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			r.log.Info("non existant user logging attempt", op, err)
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
			return
		}

		if errors.Is(err, auth.ErrInvalidCredentials) {
			r.log.Info("invalid credentials login attempt", op, err)
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
			return
		}

		r.log.Error("failed to login", op, err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to login"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"token": token})
}
