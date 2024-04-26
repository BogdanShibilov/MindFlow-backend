package consultationroute

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"

	"github.com/bogdanshibilov/mindflowbackend/internal/controller/http/v1/middleware"
	consultationservice "github.com/bogdanshibilov/mindflowbackend/internal/services/consultation"
)

type routes struct {
	log           *slog.Logger
	consultations *consultationservice.Service
}

func New(
	handler *gin.RouterGroup,
	log *slog.Logger,
	consultations *consultationservice.Service,
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
