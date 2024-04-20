package routes

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"

	"github.com/bogdanshibilov/mindflowbackend/internal/mindflow/controller/v1/dto"
	"github.com/bogdanshibilov/mindflowbackend/internal/mindflow/controller/v1/middleware"
	"github.com/bogdanshibilov/mindflowbackend/internal/mindflow/services/expert"
)

type ExpertRoutes struct {
	log     *slog.Logger
	experts *expert.Service
}

func NewExpertRoutes(handler *gin.RouterGroup, log *slog.Logger, experts *expert.Service) {
	r := &ExpertRoutes{
		log:     log,
		experts: experts,
	}

	expertsHandler := handler.Group("/expert")
	{
		expertsHandler.GET("", r.ApprovedExpertInfo)
		expertsHandler.GET("/:id", r.ExpertInfoById)
		expertsHandler.Use(middleware.RequireJwt(os.Getenv("JWTSECRET")))
		expertsHandler.POST("", r.CreateExpert)
		expertsHandler.Use(middleware.RequireAdmin())
		expertsHandler.GET("/nonapproved", r.NonApprovedExpertInfo)
		expertsHandler.POST("/approve/:id", r.ApproveExpert)
	}
}

func (r *ExpertRoutes) CreateExpert(ctx *gin.Context) {
	const op = "ExpertRoutes.CreateExpert"

	var req *dto.BecomeExpertRequest
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

	err := r.experts.CreateExpertInfo(
		ctx,
		id.(string),
		req.ChargePerHour,
		req.ExpertiseAtDescription,
	)
	if err != nil {
		r.log.Error("failed to create expert", op, err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create a new expert"})
		return
	}

	ctx.Status(http.StatusCreated)
}

func (r *ExpertRoutes) ApprovedExpertInfo(ctx *gin.Context) {
	const op = "ExpertRoutes.ApprovedExpertInfo"

	expertInfo, err := r.experts.ApprovedExpertInfo(ctx)
	if err != nil {
		r.log.Error("failed to get experts", op, err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get experts"})
		return
	}

	ctx.JSON(http.StatusOK, expertInfo)
}

func (r *ExpertRoutes) NonApprovedExpertInfo(ctx *gin.Context) {
	const op = "ExpertRoutes.NonApprovedExpertInfo"

	expertInfo, err := r.experts.NonApprovedExpertInfo(ctx)
	if err != nil {
		r.log.Error("failed to get experts", op, err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get experts"})
		return
	}

	ctx.JSON(http.StatusOK, expertInfo)
}

func (r *ExpertRoutes) ExpertInfoById(ctx *gin.Context) {
	const op = "ExpertRoutes.ExpertInfoById"

	id := ctx.Param("id")

	expert, err := r.experts.ExpertById(ctx, id)
	if err != nil {
		r.log.Error("failed to get expert", op, err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get expert"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"UserUuid":               expert.UserUuid,
		"ExpertiseAtDescription": expert.ExpertiseAtDescription,
		"ChargePerHour":          expert.ChargePerHour,
	})
}

func (r *ExpertRoutes) ApproveExpert(ctx *gin.Context) {
	const op = "ExpertRoutes.ApproveExpert"

	id := ctx.Param("id")

	err := r.experts.ApproveExpertById(ctx, id)
	if err != nil {
		r.log.Error("failed to approve expert", op, err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to approve expert"})
		return
	}

	ctx.Status(http.StatusOK)
}
