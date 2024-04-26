package expertroutes

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"

	"github.com/bogdanshibilov/mindflowbackend/internal/controller/http/v1/middleware"
	"github.com/bogdanshibilov/mindflowbackend/internal/entity"
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
		expertsHandler.GET("/filterdata", r.FilterData)
		expertsHandler.GET("/:id", r.ById)
		expertsHandler.GET("/approved", r.ExpertsWithFilter)
		expertsHandler.Use(middleware.RequireJwt(os.Getenv("JWTSECRET")))
		expertsHandler.Use(middleware.ParseClaimsIntoContext())
		expertsHandler.POST("", r.ApplyForExpert)
		expertsHandler.Use(middleware.RequireAdminPermission(users, log))
		expertsHandler.GET("", r.Experts)
		expertsHandler.PUT("/approve", r.ApproveExpert)
	}
}

func (r *routes) ById(ctx *gin.Context) {
	id := ctx.Param("id")

	expert, err := r.experts.ById(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "bad id"})
		return
	}

	dto := expertDtoFrom(expert)
	ctx.JSON(http.StatusOK, dto)
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

func (r *routes) Experts(ctx *gin.Context) {
	const op = "ExpertRoutes.Experts"

	status := getStatusQuery(ctx)
	filter := make(map[string]any)
	if status >= 0 {
		filter["status IN (?)"] = status
	}

	experts, err := r.experts.ExpertsWithFilter(ctx, filter)
	if err != nil {
		r.log.Error("failed to get pending experts", op, err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "failed to get experts"})
		return
	}

	ctx.JSON(http.StatusOK, experts)
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

func (r *routes) FilterData(ctx *gin.Context) {
	const op = "ExpertRoutes.FilterData"

	fieldsData, err := r.experts.FilterData(ctx)
	if err != nil {
		r.log.Error("failed to get filter data", op, err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "failed to get filter data"})
		return
	}

	ctx.JSON(http.StatusOK, fieldsData)
}

func (r *routes) ExpertsWithFilter(ctx *gin.Context) {
	const op = "ExpertRoutes.ExpertsWithFilter"

	filter := make(map[string]any)
	minprice := ctx.Query("minprice")
	if minprice != "" {
		filter["price >= ?"] = minprice
	}
	maxprice := ctx.Query("maxprice")
	if minprice != "" {
		filter["price <= ?"] = maxprice
	}
	filter["status IN (?)"] = entity.Approved

	experts, err := r.experts.ExpertsWithFilter(ctx, filter)
	if err != nil {
		r.log.Error("failed to get experts", op, err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "failed to get experts"})
		return
	}

	DTOs := make([]expertDTO, 0)
	for _, entity := range experts {
		DTOs = append(DTOs, *expertDtoFrom(&entity))
	}

	ctx.JSON(http.StatusOK, DTOs)
}
