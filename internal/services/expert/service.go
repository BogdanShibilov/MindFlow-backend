package expertservice

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/bogdanshibilov/mindflowbackend/internal/entity"
	expertrepo "github.com/bogdanshibilov/mindflowbackend/internal/repository/expert"
	userrepo "github.com/bogdanshibilov/mindflowbackend/internal/repository/user"
)

type Service struct {
	expertRepo *expertrepo.Repo
	userRepo   *userrepo.Repo
}

func New(expertRepo *expertrepo.Repo, userRepo *userrepo.Repo) *Service {
	return &Service{
		expertRepo: expertRepo,
		userRepo:   userRepo,
	}
}

func (s *Service) ApplyForExpert(
	ctx context.Context,
	userId string,
	helpDescription string,
	price int,
) error {
	const op = "services.expert.ApplyForExpert"

	uuid, err := uuid.Parse(userId)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	information := &entity.ExpertInformation{
		Price:           price,
		HelpDescription: helpDescription,
	}
	application := &entity.ExpertApplication{
		Status: entity.Pending,
	}
	expert := &entity.Expert{
		UserUuid:          uuid,
		ExpertInformation: *information,
		ExpertApplication: *application,
	}

	err = s.expertRepo.CreateExpert(ctx, expert)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Service) ChangeExpertStatus(ctx context.Context, expertId string, status entity.Status) error {
	const op = "services.expert.ApproveExpert"

	application := &entity.ExpertApplication{
		Status: status,
	}

	uuid, err := uuid.Parse(expertId)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return s.expertRepo.UpdateExpertApplication(ctx, application, uuid)
}

func (s *Service) ById(ctx context.Context, expertId string) (*entity.Expert, error) {
	const op = "services.expert.ApproveExpert"

	uuid, err := uuid.Parse(expertId)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return s.expertRepo.ByUuid(ctx, uuid)
}

func (s *Service) FilterData(ctx context.Context) (*FilterData, error) {
	const op = "services.expert.FilterData"

	profFields, err := s.expertRepo.ProffFieldListAndCount(ctx)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	minMaxPrice, err := s.expertRepo.MinMaxPrice(ctx)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &FilterData{
		ProffFieldData: profFields,
		MinMaxPrice:    *minMaxPrice,
	}, nil
}

func (s *Service) ExpertsWithFilter(ctx context.Context, filter map[string]any) ([]entity.Expert, error) {
	return s.expertRepo.ExpertsWithFilter(ctx, filter)
}

func (s *Service) DoesExist(ctx context.Context, id string) (bool, error) {
	const op = "services.expert.DoesExist"

	uuid, err := uuid.Parse(id)
	if err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}

	return s.expertRepo.DoesExist(ctx, uuid)
}
