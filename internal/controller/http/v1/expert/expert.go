package expertroutes

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"

	"github.com/bogdanshibilov/mindflowbackend/internal/controller/http/v1/middleware"
	"github.com/bogdanshibilov/mindflowbackend/internal/entity"
	expertrepo "github.com/bogdanshibilov/mindflowbackend/internal/repository/expert"
	expertservice "github.com/bogdanshibilov/mindflowbackend/internal/services/expert"
	userservice "github.com/bogdanshibilov/mindflowbackend/internal/services/user"
)

type routes struct {
	log     *slog.Logger
	experts *expertservice.Service
}

func New(
	handler *gin.RouterGroup,
	log *slog.Logger,
	experts *expertservice.Service,
	users *userservice.Service,
) {
	r := &routes{
		log:     log,
		experts: experts,
	}

	expertsHandler := handler.Group("/experts")
	{
		expertsHandler.GET("/approved", r.ApprovedExperts)
		expertsHandler.Use(middleware.RequireJwt(os.Getenv("JWTSECRET")))
		expertsHandler.Use(middleware.ParseClaimsIntoContext())
		expertsHandler.POST("", r.ApplyForExpert)
		expertsHandler.Use(middleware.RequireAdminPermission(users, log))
		expertsHandler.GET("", r.Experts)
		expertsHandler.PUT("/approve", r.ApproveExpert)
	}
}

func (r *routes) ApplyForExpert(ctx *gin.Context) {
	const op = "ExpertRoutes.ApplyForExpert"

	var req *applyForExpertRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		r.log.Warn("invalid JSON received", op, err)
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "invalid JSON"})
		return
	}

	id := ctx.GetString("uuid")

	err := r.experts.ApplyForExpert(
		ctx,
		id,
		req.HelpDescription,
		req.Price,
	)

	if err != nil {
		r.log.Error("failed to appy for expert", op, err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "failed to appy for expert"})
		return
	}

	ctx.Status(http.StatusCreated)
}

func (r *routes) ApprovedExperts(ctx *gin.Context) {
	const op = "ExpertRoutes.AcceptedExperts"

	experts, err := r.experts.Experts(
		ctx,
		expertrepo.SelectStatus(entity.Approved),
	)
	if err != nil {
		r.log.Error("failed to get approved experts", op, err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "failed to get experts"})
		return
	}

	DTOs := make([]expertDTO, 0)
	for _, entity := range experts {
		DTOs = append(DTOs, *expertDtoFrom(&entity))
	}

	ctx.JSON(http.StatusOK, DTOs)
}

func (r *routes) Experts(ctx *gin.Context) {
	const op = "ExpertRoutes.Experts"

	status, err := getStatusQuery(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "invalid status query"})
		return
	}

	experts, err := r.experts.Experts(
		ctx,
		expertrepo.SelectStatus(status),
	)
	if err != nil {
		r.log.Error("failed to get pending experts", op, err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "failed to get experts"})
		return
	}

	DTOs := make([]expertDTO, 0)
	for _, entity := range experts {
		DTOs = append(DTOs, *expertDtoFrom(&entity))
	}

	ctx.JSON(http.StatusOK, DTOs)
}

func (r *routes) ApproveExpert(ctx *gin.Context) {
	const op = "ExpertRoutes.ApproveExpert"

	var req *approveExpertRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		r.log.Warn("invalid JSON received", op, err)
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "invalid JSON"})
		return
	}

	err := r.experts.ApproveExpert(ctx, req.ExpertId)
	if err != nil {
		r.log.Error("failed to approve expert", op, err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "failed to approve expert"})
		return
	}

	ctx.Status(http.StatusOK)
}