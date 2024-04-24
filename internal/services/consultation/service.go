package consultationservice

import (
	"context"
	"fmt"

	"github.com/bogdanshibilov/mindflowbackend/internal/entity"
	consultationrepo "github.com/bogdanshibilov/mindflowbackend/internal/repository/consultation"
	"github.com/google/uuid"
)

type Service struct {
	consultRepo *consultationrepo.Repo
}

func New(consultRepo consultationrepo.Repo) *Service {
	return &Service{
		consultRepo: &consultRepo,
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

func (s *Service) ByPersonId(ctx context.Context, id string, opts ...consultationrepo.ByPersonUuidOption) ([]entity.Consultation, error) {
	const op = "services.consultation.ById"

	uuid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return s.consultRepo.ByPersonUuid(ctx, uuid, opts...)
}

func (s *Service) ApproveApplication(ctx context.Context, id string) error {
	const op = "services.consultation.ApproveApplication"

	uuid, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return s.consultRepo.UpdateApplicationStatus(ctx, uuid, entity.Approved)
}
