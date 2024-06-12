package consultationroute

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"

	"github.com/bogdanshibilov/mindflowbackend/internal/controller/http/v1/middleware"
	"github.com/bogdanshibilov/mindflowbackend/internal/entity"
	consultationrepo "github.com/bogdanshibilov/mindflowbackend/internal/repository/consultation"
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
		consultHandler.GET("alreadyapplied/:expertid", r.AlreadyApplied)
		consultHandler.GET("meetasstudent", r.MeetingsAsStudent)
		consultHandler.GET("meetasexpert", r.MeetingsAsExpert)
		consultHandler.GET("/:id", r.ById)
		consultHandler.Use(middleware.RequireAdminPermission(userservice, log))
		consultHandler.GET("", r.Consultations)
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

func (r *routes) MeetingsAsStudent(ctx *gin.Context) {
	id := ctx.GetString("uuid")

	consults, err := r.consultations.ByPersonId(
		ctx,
		id,
		consultationrepo.SelectByWhoseUuid(consultationrepo.ByMenteeUuid),
	)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "bad request"})
		return
	}

	var meetings []entity.ConsultationMeeting
	for _, consult := range consults {
		meets, err := r.consultations.MeetingsByConsultationId(ctx, consult.Uuid.String())
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": "bad request"})
			return
		}
		meetings = append(meetings, meets...)
	}

	ctx.JSON(http.StatusOK, meetings)
}

func (r *routes) MeetingsAsExpert(ctx *gin.Context) {
	id := ctx.GetString("uuid")

	consults, err := r.consultations.ByPersonId(
		ctx,
		id,
		consultationrepo.SelectByWhoseUuid(consultationrepo.ByExpertUuid),
	)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "bad request"})
		return
	}

	var meetings []entity.ConsultationMeeting
	for _, consult := range consults {
		meets, err := r.consultations.MeetingsByConsultationId(ctx, consult.Uuid.String())
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": "bad request"})
			return
		}
		meetings = append(meetings, meets...)
	}

	ctx.JSON(http.StatusOK, meetings)
}

func (r *routes) AlreadyApplied(ctx *gin.Context) {
	const op = "consultationroutes.AlreadyApplied"

	menteeId := ctx.GetString("uuid")
	expertId := ctx.Param("expertid")

	alreadyApplied, err := r.consultations.DoesExist(ctx, menteeId, expertId)
	if err != nil {
		r.log.Error("failed to check existance", op, err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "failed to check existance"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"alreadyApplied": alreadyApplied})
}
