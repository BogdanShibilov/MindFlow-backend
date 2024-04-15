package routes

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/bogdanshibilov/mindflowbackend/internal/mindflow/controller/v1/dto"
	"github.com/bogdanshibilov/mindflowbackend/internal/mindflow/controller/v1/middleware"
	"github.com/bogdanshibilov/mindflowbackend/internal/mindflow/services/auth"
	"github.com/gin-gonic/gin"
)

type UserRoutes struct {
	log  *slog.Logger
	auth *auth.Auth
}

func NewUserRoutes(handler *gin.RouterGroup, log *slog.Logger, auth *auth.Auth) {
	r := &UserRoutes{
		log:  log,
		auth: auth,
	}

	userHandler := handler.Group("/user")
	{
		userHandler.Use(middleware.RequireJwt(os.Getenv("JWTSECRET")))
		userHandler.GET("/:id", r.UserDetailsById)
		userHandler.POST("/updatedetails", r.UpdateUserDetails)
	}
}

func (r *UserRoutes) UpdateUserDetails(ctx *gin.Context) {
	const op = "AuthRoutes.UpdateUserDetails"

	var req *dto.UpdateUserDetailsRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		r.log.Error("failed to bind json data", op, err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON"})
		return
	}

	err := r.auth.UpdateUserDetails(
		ctx,
		req.UserId,
		req.PhoneNumber,
		req.ProfessionalField,
		req.ExperienceDescription,
	)
	if err != nil {
		r.log.Error("failed to update user details", op, err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update user details"})
		return
	}

	ctx.Status(http.StatusOK)
}

func (r *UserRoutes) UserDetailsById(ctx *gin.Context) {
	const op = "AuthRoutes.UserDetailsById"

	id := ctx.Param("id")

	userDetails, err := r.auth.UserDetailsByUserUuid(ctx, id)
	if err != nil {
		r.log.Error("failed to get expert", op, err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get expert"})
		return
	}

	ctx.JSON(http.StatusOK, userDetails)
}
