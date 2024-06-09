package consultationservice

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/bogdanshibilov/mindflowbackend/internal/entity"
	consultationrepo "github.com/bogdanshibilov/mindflowbackend/internal/repository/consultation"
	userrepo "github.com/bogdanshibilov/mindflowbackend/internal/repository/user"
	"github.com/bogdanshibilov/mindflowbackend/internal/services/mails"
)

type Service struct {
	consultRepo *consultationrepo.Repo
	userRepo    *userrepo.Repo
}

func New(consultRepo consultationrepo.Repo, userRepo userrepo.Repo) *Service {
	return &Service{
		consultRepo: &consultRepo,
		userRepo:    &userRepo,
	}
}

func (s *Service) ApplyForConsultation(
	ctx context.Context,
	menteeId string,
	expertId string,
	menteeQuestions string,
) error {
	const op = "services.consultation.ApplyForConsultation"

	menteeUuid, err := uuid.Parse(menteeId)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	expertUuid, err := uuid.Parse(expertId)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	application := entity.ConsultationApplication{
		Status:          entity.Pending,
		MenteeQuestions: menteeQuestions,
	}
	consultation := entity.Consultation{
		ExpertUuid:              expertUuid,
		MenteeUuid:              menteeUuid,
		ConsultationApplication: application,
	}

	err = s.consultRepo.CreateConsultation(ctx, &consultation)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Service) ById(ctx context.Context, id string) (*entity.Consultation, error) {
	const op = "services.consultation.ById"

	uuid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return s.consultRepo.ByUuid(ctx, uuid)
}

func (s *Service) Consultations(ctx context.Context) ([]entity.Consultation, error) {
	return s.consultRepo.Consultations(ctx)
}

func (s *Service) ByPersonId(ctx context.Context, id string, opts ...consultationrepo.ByPersonUuidOption) ([]entity.Consultation, error) {
	const op = "services.consultation.ById"

	uuid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return s.consultRepo.ByPersonUuid(ctx, uuid, opts...)
}

func (s *Service) ChangeApplicationStatus(ctx context.Context, id string, status entity.Status) error {
	const op = "services.consultation.ChangeApplicationStatus"

	uuid, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return s.consultRepo.UpdateApplicationStatus(ctx, uuid, status)
}

func (s *Service) CreateMeeting(
	ctx context.Context,
	consultId string,
	startTime time.Time,
	link string,
) error {
	const op = "services.consultation.CreateMeeting"

	uuid, err := uuid.Parse(consultId)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	meeting := &entity.ConsultationMeeting{
		ConsultationUuid: uuid,
		StartTime:        startTime,
		Link:             link,
	}

	err = s.consultRepo.CreateMeeting(ctx, meeting)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	err = s.ChangeApplicationStatus(ctx, consultId, entity.Approved)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	consult, err := s.consultRepo.ByUuid(ctx, uuid)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	expert, err := s.userRepo.ByUuid(ctx, consult.ExpertUuid)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	mentee, err := s.userRepo.ByUuid(ctx, consult.MenteeUuid)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	mails.SendConsultationNotification([]string{expert.Email, mentee.Email}, meeting.Link)

	return nil
}

func (s *Service) MeetingsByConsultationId(ctx context.Context, id string) ([]entity.ConsultationMeeting, error) {
	const op = "services.consultation.MeetingsByConsultationId"

	uuid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return s.consultRepo.MeetingsByConsultationUuid(ctx, uuid)
}
