package routes

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/bogdanshibilov/mindflowbackend/internal/mindflow/controller/v1/dto"
	"github.com/bogdanshibilov/mindflowbackend/internal/mindflow/controller/v1/middleware"
	"github.com/bogdanshibilov/mindflowbackend/internal/mindflow/entity"
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
		enrollmentsHandler.POST("/", r.Enroll)
		enrollmentsHandler.GET("/mentor/meetings", r.MeetingsAsMentor)
		enrollmentsHandler.GET("/mentee/meetings", r.MeetingsAsMentee)
		enrollmentsHandler.Use(middleware.RequireAdmin())
		enrollmentsHandler.GET("/", r.EnrollmentsByMemberId)
		enrollmentsHandler.POST("/approve/:id", r.ApproveEnrollment)
	}
}

func (r *EnrollmentRoutes) Enroll(ctx *gin.Context) {
	const op = "ExpertRoutes.Enroll"

	var req *dto.CreateEnrollmentRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		r.log.Error("failed to bind json data", op, err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON"})
		return
	}

	id, ok := ctx.Get("userId")
	if !ok {
		r.log.Error("failed to get user id", op, "no userId in context")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get user id"})
		return
	}

	err := r.enrollments.CreateEnrollment(
		ctx,
		req.MentorId,
		id.(string),
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

func (r *EnrollmentRoutes) ApproveEnrollment(ctx *gin.Context) {
	const op = "ExpertRoutes.ApproveEnrollment"

	id := ctx.Param("id")

	err := r.enrollments.ApproveEnrollmentById(ctx, id)
	if err != nil {
		r.log.Error("failed to approve enrollment", op, err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to approve enrollment"})
		return
	}

	ctx.Status(http.StatusOK)
}

func (r *EnrollmentRoutes) MeetingsAsMentor(ctx *gin.Context) {
	const op = "ExpertRoutes.ApproveEnrollment"

	roles, ok := ctx.Get("roles")
	if !ok {
		r.log.Error("failed to get user roles", op, "no roles in context")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get user roles"})
		return
	}

	if !contains(roles.([]string), "expert") {
		r.log.Error("user is not expert", op, "user is not expert")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "user is not expert"})
		return
	}

	id, ok := ctx.Get("userId")
	if !ok {
		r.log.Error("failed to get user id", op, "no userId in context")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get user id"})
		return
	}

	enrollmentsAsMentor, err := r.enrollments.EnrollmentsByMemberId(ctx, id.(string), enrollment.Mentor)
	if err != nil {
		r.log.Error("failed to get enrollments", op, "failed to get enrollments")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get enrollments"})
		return
	}

	var meetings []entity.Meeting
	for _, enrollment := range enrollmentsAsMentor {
		m, err := r.enrollments.MeetingsByEnrollmentUuid(ctx, enrollment.Uuid.String())
		if err != nil {
			r.log.Error("failed to get meetings", op, "failed to get meetings")
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get meetings"})
			return
		}
		meetings = append(meetings, m...)
	}

	ctx.JSON(http.StatusOK, meetings)
}

func (r *EnrollmentRoutes) MeetingsAsMentee(ctx *gin.Context) {
	const op = "ExpertRoutes.ApproveEnrollment"

	id, ok := ctx.Get("userId")
	if !ok {
		r.log.Error("failed to get user id", op, "no userId in context")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get user id"})
		return
	}

	enrollmentsAsMentee, err := r.enrollments.EnrollmentsByMemberId(ctx, id.(string), enrollment.Mentee)
	if err != nil {
		r.log.Error("failed to get enrollments", op, "failed to get enrollments")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get enrollments"})
		return
	}

	var meetings []entity.Meeting
	for _, enrollment := range enrollmentsAsMentee {
		m, err := r.enrollments.MeetingsByEnrollmentUuid(ctx, enrollment.Uuid.String())
		if err != nil {
			r.log.Error("failed to get meetings", op, "failed to get meetings")
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get meetings"})
			return
		}
		meetings = append(meetings, m...)
	}

	ctx.JSON(http.StatusOK, meetings)
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
