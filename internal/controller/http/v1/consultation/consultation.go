package consultationroute

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"

	"github.com/bogdanshibilov/mindflowbackend/internal/controller/http/v1/middleware"
	"github.com/bogdanshibilov/mindflowbackend/internal/entity"
	consultationservice "github.com/bogdanshibilov/mindflowbackend/internal/services/consultation"
	userservice "github.com/bogdanshibilov/mindflowbackend/internal/services/user"
)

type routes struct {
	log           *slog.Logger
	consultations *consultationservice.Service
}

func New(
	handler *gin.RouterGroup,
	log *slog.Logger,
	consultations *consultationservice.Service,
	userservice *userservice.Service,
) {
	r := &routes{
		log:           log,
		consultations: consultations,
	}

	consultHandler := handler.Group("/consultation")
	{
		consultHandler.Use(middleware.RequireJwt(os.Getenv("JWTSECRET")))
		consultHandler.Use(middleware.ParseClaimsIntoContext())
		consultHandler.POST("apply", r.ApplyForConsultation)
		consultHandler.Use(middleware.RequireAdminPermission(userservice, log))
		consultHandler.GET("", r.Consultations)
		consultHandler.GET("/:id", r.ById)
		consultHandler.POST("/meeting", r.CreateMeeting)
		consultHandler.POST("/reject/:id", r.RejectApplication)
	}
}

func (r *routes) ApplyForConsultation(ctx *gin.Context) {
	const op = "consultationroutes.ApplyForConsultation"

	var req *applyForConsultationRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		r.log.Warn("invalid JSON received", op, err)
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "invalid JSON"})
		return
	}

	id := ctx.GetString("uuid")

	err := r.consultations.ApplyForConsultation(ctx, id, req.ExpertId, req.MenteeQuestions)
	if err != nil {
		r.log.Error("failed to appy for consultation", op, err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "failed to appy for consultation"})
		return
	}

	ctx.Status(http.StatusCreated)
}

func (r *routes) Consultations(ctx *gin.Context) {
	const op = "consultationroutes.ApplyForConsultation"

	consults, err := r.consultations.Consultations(ctx)
	if err != nil {
		r.log.Error("failed to get consultations", op, err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "failed to get consultations"})
		return
	}

	ctx.JSON(http.StatusOK, consults)
}

func (r *routes) ById(ctx *gin.Context) {
	id := ctx.Param("id")

	consult, err := r.consultations.ById(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "bad id"})
		return
	}

	ctx.JSON(http.StatusOK, consult)
}

func (r *routes) CreateMeeting(ctx *gin.Context) {
	const op = "consultationroutes.CreateMeeting"

	var req *createMeetingRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		r.log.Warn("invalid JSON received", op, err)
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "invalid JSON"})
		return
	}

	err := r.consultations.CreateMeeting(
		ctx,
		req.ConsultationId,
		req.StartTime,
		req.Link,
	)
	if err != nil {
		r.log.Error("failed to create a meeting", op, err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "failed to create a meeting"})
		return
	}

	ctx.Status(http.StatusOK)
}

func (r *routes) RejectApplication(ctx *gin.Context) {
	id := ctx.Param("id")

	err := r.consultations.ChangeApplicationStatus(ctx, id, entity.Rejected)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "bad id"})
		return
	}

	ctx.Status(http.StatusOK)
}
