package userroutes

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"

	"github.com/bogdanshibilov/mindflowbackend/internal/controller/http/v1/middleware"
	userservice "github.com/bogdanshibilov/mindflowbackend/internal/services/user"
)

type routes struct {
	log   *slog.Logger
	users *userservice.Service
}

func New(
	handler *gin.RouterGroup,
	log *slog.Logger,
	users *userservice.Service,
) {
	r := &routes{
		log:   log,
		users: users,
	}

	usersHandler := handler.Group("/users")
	{
		usersHandler.Use(middleware.RequireJwt(os.Getenv("JWTSECRET")))
		usersHandler.GET("/:id", r.ById)
		usersHandler.Use(middleware.RequireAdminPermission(users, log))
		usersHandler.GET("", r.Users)
		usersHandler.PUT("/forceupdateuserprofile", r.ForceUpdateUserProfile)
		usersHandler.DELETE("", r.DeleteUserById)
	}
}

func (r *routes) Users(ctx *gin.Context) {
	const op = "UserRoutes.Users"

	users, err := r.users.Users(ctx)
	if err != nil {
		r.log.Error("failed to get users", op, err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "failed to get users"})
		return
	}

	DTOs := make([]userDto, 0)
	for _, entity := range users {
		DTOs = append(DTOs, *userDtoFrom(&entity))
	}

	ctx.JSON(http.StatusOK, DTOs)
}

func (r *routes) ForceUpdateUserProfile(ctx *gin.Context) {
	const op = "UserRoutes.ForceUpdateUserProfile"

	var req *ForceUpdateUserProfileRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		r.log.Warn("invalid JSON received", op, err)
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "invalid JSON"})
		return
	}

	err := r.users.UpdateProfile(
		ctx,
		req.Name,
		req.Email,
		req.Phone,
		req.ProfessionalField,
		req.ExperienceDescription,
		req.Id,
	)
	if err != nil {
		r.log.Error("failed to update user", op, err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "failed to update user"})
		return
	}

	ctx.Status(http.StatusOK)
}

func (r *routes) DeleteUserById(ctx *gin.Context) {
	const op = "UserRoutes.DeleteUserById"

	var req *DeleteUserByIdRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		r.log.Warn("invalid JSON received", op, err)
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "invalid JSON"})
		return
	}

	err := r.users.DeleteUserById(ctx, req.Id)
	if err != nil {
		r.log.Error("failed to delete user", op, err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "failed to delete user"})
		return
	}

	ctx.Status(http.StatusOK)
}

func (r *routes) ById(ctx *gin.Context) {
	id := ctx.Param("id")

	user, err := r.users.ById(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "bad id"})
		return
	}

	ctx.JSON(http.StatusOK, userDtoFrom(user))
}
