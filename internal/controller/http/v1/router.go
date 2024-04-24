package v1

import (
	"log/slog"
	"net/http"

	"github.com/gin-contrib/cors"

	"github.com/gin-gonic/gin"

	authroutes "github.com/bogdanshibilov/mindflowbackend/internal/controller/http/v1/auth"
	expertroutes "github.com/bogdanshibilov/mindflowbackend/internal/controller/http/v1/expert"
	authservice "github.com/bogdanshibilov/mindflowbackend/internal/services/auth"
	expertservice "github.com/bogdanshibilov/mindflowbackend/internal/services/expert"
	userservice "github.com/bogdanshibilov/mindflowbackend/internal/services/user"
)

func NewRouter(
	handler *gin.Engine,
	log *slog.Logger,
	auth *authservice.Service,
	experts *expertservice.Service,
	users *userservice.Service,
) {
	handler.Use(gin.Recovery())

	handler.Use(cors.New(cors.Config{
		AllowWildcard: true,
		AllowOrigins:  []string{"http://localhost:*"},
		AllowMethods:  []string{"GET", "POST", "PUT"},
		AllowHeaders:  []string{"Origin", "authorization", "content-type", "accept"},
	}))

	handler.GET("/healthz", func(ctx *gin.Context) { ctx.Status(http.StatusOK) })

	h := handler.Group("/api/v1")
	{
		authroutes.New(h, log, auth)
		expertroutes.New(h, log, experts, users)
	}
}
