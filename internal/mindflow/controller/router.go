package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func New(handler *gin.Engine) {
	handler.Use(gin.Recovery())

	handler.GET("/healthz", func(ctx *gin.Context) { ctx.Status(http.StatusOK) })

	handler.GET("/hello", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"message": "Hello!"})
	})
}
