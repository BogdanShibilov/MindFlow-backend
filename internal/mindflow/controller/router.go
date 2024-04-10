package controller

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/bogdanshibilov/mindflowbackend/internal/mindflow/controller/v1/routes"
	"github.com/bogdanshibilov/mindflowbackend/internal/mindflow/services/auth"
	"github.com/bogdanshibilov/mindflowbackend/internal/mindflow/services/expert"
)

func New(
	handler *gin.Engine,
	log *slog.Logger,
	auth *auth.Auth,
	experts *expert.Service,
) {
	handler.Use(gin.Recovery())

	handler.GET("/healthz", func(ctx *gin.Context) { ctx.Status(http.StatusOK) })

	h := handler.Group("/api/v1")
	{
		routes.NewAuthRoutes(h, log, auth)
		routes.NewExpertRoutes(h, log, experts)
	}
}
