package routes

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/bogdanshibilov/mindflowbackend/internal/mindflow/controller/v1/dto"
	"github.com/bogdanshibilov/mindflowbackend/internal/mindflow/controller/v1/middleware"
	"github.com/bogdanshibilov/mindflowbackend/internal/mindflow/services/enrollment"
	"github.com/gin-gonic/gin"
)

type EnrollmentRoutes struct {
	log         *slog.Logger
	enrollments *enrollment.Service
}

func NewEnrollmentRoutes(handler *gin.RouterGroup, log *slog.Logger, enrollments *enrollment.Service) {
	r := &EnrollmentRoutes{
		log:         log,
		enrollments: enrollments,
	}

	enrollmentsHandler := handler.Group("/enrollment")
	{
		enrollmentsHandler.Use(middleware.RequireJwt(os.Getenv("JWTSECRET")))
		enrollmentsHandler.POST("/", r.CreateEnrollment)
		enrollmentsHandler.GET("/", r.EnrollmentsByMemberId)
	}
}

func (r *EnrollmentRoutes) CreateEnrollment(ctx *gin.Context) {
	const op = "ExpertRoutes.CreateEnrollment"

	var req *dto.CreateEnrollmentRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		r.log.Error("failed to bind json data", op, err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON"})
		return
	}

	err := r.enrollments.CreateEnrollment(
		ctx,
		req.MentorId,
		req.MenteeId,
		req.MenteeQuestions,
	)
	if err != nil {
		r.log.Error("failed to create expert", op, err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create a new expert"})
		return
	}

	ctx.Status(http.StatusCreated)
}

func (r *EnrollmentRoutes) EnrollmentsByMemberId(ctx *gin.Context) {
	const op = "ExpertRoutes.EnrollmentsByMemberId"

	var req *dto.EnrollmentsByMemberIdRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		r.log.Error("failed to bind json data", op, err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON"})
		return
	}

	if req.ByWhoseId != string(enrollment.Mentor) &&
		req.ByWhoseId != string(enrollment.Mentee) {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid ByWhoseId parameter "})
		return
	}

	enrollments, err := r.enrollments.EnrollmentsByMemberId(ctx, req.Id, enrollment.ByWho(req.ByWhoseId))
	if err != nil {
		r.log.Error("failed to get enrollments", op, err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get enrollments"})
		return
	}

	ctx.JSON(http.StatusOK, enrollments)
}
