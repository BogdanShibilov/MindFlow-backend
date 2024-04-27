package userroutes

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/bogdanshibilov/mindflowbackend/internal/controller/http/v1/middleware"
	userservice "github.com/bogdanshibilov/mindflowbackend/internal/services/user"
	"github.com/gin-gonic/gin"
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
		usersHandler.Use(middleware.RequireAdminPermission(users, log))
		usersHandler.GET("", r.Users)
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
