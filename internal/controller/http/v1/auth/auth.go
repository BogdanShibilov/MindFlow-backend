package authroutes

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"

	userrepo "github.com/bogdanshibilov/mindflowbackend/internal/repository/user"
	authservice "github.com/bogdanshibilov/mindflowbackend/internal/services/auth"
)

type routes struct {
	log  *slog.Logger
	auth *authservice.Service
}

func New(
	handler *gin.RouterGroup,
	log *slog.Logger,
	auth *authservice.Service,
) {
	r := &routes{
		log:  log,
		auth: auth,
	}

	authHandler := handler.Group("/auth")
	{
		authHandler.POST("/signup", r.SignUp)
		authHandler.POST("/emailsignin", r.SignInWithEmail)
	}
}

func (r *routes) SignUp(ctx *gin.Context) {
	const op = "AuthRoutes.SignUp"

	var req *signUpRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		r.log.Warn("invalid JSON received", op, err)
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "invalid JSON"})
		return
	}

	err := r.auth.Register(
		ctx,
		req.Username,
		req.Password,
		req.Email,
		req.Phone,
		req.ProfessionalField,
		req.ExperienceDescription,
	)
	if err != nil {
		r.log.Error("failed to register a new user", op, err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "failed to register"})
		return
	}

	ctx.Status(http.StatusCreated)
}

func (r *routes) SignInWithEmail(ctx *gin.Context) {
	const op = "AuthRoutes.SignIn"

	var req *signInWithEmailRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		r.log.Warn("invalid JSON received", op, err)
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "invalid JSON"})
		return
	}

	token, err := r.auth.LoginByEmail(ctx, req.Email, req.Password)
	if err != nil {
		if errors.Is(err, authservice.ErrInvalidCredentials) {
			r.log.Warn("invalid credentials received", op, err)
			ctx.JSON(http.StatusUnauthorized, gin.H{"message": "invalid credentials"})
			return
		} else if errors.Is(err, userrepo.ErrUserNotFound) {
			r.log.Warn("nonexisting credentials received", op, err)
			ctx.JSON(http.StatusUnauthorized, gin.H{"message": "invalid credentials"})
			return
		} else {
			r.log.Error("failed to login a new user", op, err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": "failed to sign in"})
			return
		}
	}

	ctx.JSON(http.StatusOK, signInWithEmailResponse{
		AccessToken: token,
	})
}
