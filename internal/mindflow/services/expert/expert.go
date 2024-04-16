package expert

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/bogdanshibilov/mindflowbackend/internal/mindflow/entity"
)

type Service struct {
	expertSaver    ExpertSaver
	expertProvider ExpertProvider
	users          Users
}

type ExpertSaver interface {
	SaveExpertInfo(ctx context.Context, expertInfo *entity.ExpertInfo) error
	UpdateExpertInfo(ctx context.Context, info *entity.ExpertInfo) error
}

type ExpertProvider interface {
	ExpertByUuid(ctx context.Context, uuid *uuid.UUID) (*entity.ExpertInfo, error)
	ExpertInfo(ctx context.Context) ([]entity.ExpertInfo, error)
}

type Users interface {
	UpdateUser(ctx context.Context, user *entity.User) error
	UserByUuid(ctx context.Context, uuid uuid.UUID) (*entity.User, error)
}

func New(
	expertSaver ExpertSaver,
	expertProvider ExpertProvider,
	userUpdater Users,
) *Service {
	return &Service{
		expertSaver:    expertSaver,
		expertProvider: expertProvider,
		users:          userUpdater,
	}
}

func (s *Service) CreateExpertInfo(
	ctx context.Context,
	expertId string,
	chargePerHour int,
	expertiseAtDesc string,
) error {
	const op = "Expert.CreateExpertInfo"

	uuid, err := uuid.Parse(expertId)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	expertInfo := &entity.ExpertInfo{
		UserUuid:               uuid,
		ChargePerHour:          chargePerHour,
		ExpertiseAtDescription: expertiseAtDesc,
	}

	err = s.expertSaver.SaveExpertInfo(ctx, expertInfo)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Service) ExpertInfo(ctx context.Context) ([]entity.ExpertInfo, error) {
	const op = "Expert.ExpertInfo"

	expertInfo, err := s.expertProvider.ExpertInfo(ctx)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return expertInfo, nil
}

func (s *Service) ExpertById(ctx context.Context, id string) (*entity.ExpertInfo, error) {
	const op = "Expert.ExpertById"

	uuid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	expertInfo, err := s.expertProvider.ExpertByUuid(ctx, &uuid)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return expertInfo, nil
}

func (s *Service) ApproveExpertById(ctx context.Context, id string) error {
	const op = "Expert.ApproveExpertById"

	uuid, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	expertInfo, err := s.expertProvider.ExpertByUuid(ctx, &uuid)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	expertInfo.IsApproved = true
	err = s.expertSaver.UpdateExpertInfo(ctx, expertInfo)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	user, err := s.users.UserByUuid(ctx, uuid)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	user.Roles = append(user.Roles, "expert")
	err = s.users.UpdateUser(ctx, user)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
