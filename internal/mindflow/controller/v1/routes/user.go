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
		userHandler.GET("/:id", r.UserByIdNoJwt)
		userHandler.GET("/userdetails/:id", r.UserDetailsByIdNoJwt)
		userHandler.Use(middleware.RequireJwt(os.Getenv("JWTSECRET")))
		userHandler.GET("/userdetails", r.UserDetailsById)
		userHandler.POST("/updatemydetails", r.UpdateUserDetails)
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

	id, ok := ctx.Get("userId")
	if !ok {
		r.log.Error("failed to get user details", op, "no userId in context")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get user details"})
		return
	}

	err := r.auth.UpdateUserDetails(
		ctx,
		id.(string),
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

	id, ok := ctx.Get("userId")
	if !ok {
		r.log.Error("failed to get user details", op, "no userId in context")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get user details"})
		return
	}

	userDetails, err := r.auth.UserDetailsByUserId(ctx, id.(string))
	if err != nil {
		r.log.Error("failed to get user details", op, err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get user details"})
		return
	}

	ctx.JSON(http.StatusOK, userDetails)
}

func (r *UserRoutes) UserByIdNoJwt(ctx *gin.Context) {
	const op = "AuthRoutes.UserByIdNoJwt"

	id := ctx.Param("id")

	user, err := r.auth.UserById(ctx, id)
	if err != nil {
		r.log.Error("failed to get user", op, err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get user"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"id":    user.Uuid,
		"email": user.Email,
		"name":  user.Name,
	})
}

func (r *UserRoutes) UserDetailsByIdNoJwt(ctx *gin.Context) {
	const op = "AuthRoutes.UserDetailsByIdNoJwt"

	id := ctx.Param("id")

	userDetails, err := r.auth.UserDetailsByUserId(ctx, id)
	if err != nil {
		r.log.Error("failed to get user details", op, err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get user details"})
		return
	}

	ctx.JSON(http.StatusOK, userDetails)
}
