package userservice

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/bogdanshibilov/mindflowbackend/internal/entity"
	userrepo "github.com/bogdanshibilov/mindflowbackend/internal/repository/user"
)

type Service struct {
	userRepo *userrepo.Repo
}

func New(userRepo *userrepo.Repo) *Service {
	return &Service{
		userRepo: userRepo,
	}
}

func (s *Service) CreateUser(
	ctx context.Context,
	username string,
	passHash string,
	name string,
	email string,
	phone string,
	professionalField string,
	experienceDescription string,
) error {
	const op = "services.user.CreateUser"

	userCreds := entity.UserCredentials{
		Username: username,
		PassHash: []byte(passHash),
	}
	userProfile := entity.UserProfile{
		Name:                  name,
		Email:                 email,
		Phone:                 phone,
		ProfessionalField:     professionalField,
		ExperienceDescription: experienceDescription,
	}
	newUser := &entity.User{
		UserCredentials: userCreds,
		UserProfile:     userProfile,
	}

	err := s.userRepo.CreateUser(ctx, newUser)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (s *Service) UpdateCredentials(
	ctx context.Context,
	newUsername string,
	newPassHash string,
	userId string,
) error {
	const op = "services.user.UpdateCredentials"

	uuid, err := uuid.Parse(userId)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	newCreds := &entity.UserCredentials{
		Username: newUsername,
		PassHash: []byte(newPassHash),
	}

	err = s.userRepo.UpdateCredentials(ctx, newCreds, uuid)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Service) UpdateProfile(
	ctx context.Context,
	newName string,
	newEmail string,
	newPhone string,
	newProfessionalField string,
	newExperienceDescription string,
	userId string,
) error {
	const op = "services.user.UpdateProfile"

	uuid, err := uuid.Parse(userId)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	newProfile := &entity.UserProfile{
		Name:                  newName,
		Email:                 newEmail,
		Phone:                 newPhone,
		ProfessionalField:     newProfessionalField,
		ExperienceDescription: newExperienceDescription,
	}

	err = s.userRepo.UpdateProfile(ctx, newProfile, uuid)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Service) ById(ctx context.Context, id string) (*entity.User, error) {
	const op = "services.user.ById"

	uuid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	user, err := s.userRepo.ByUuid(ctx, uuid)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil
}

func (s *Service) ByEmail(ctx context.Context, email string) (*entity.User, error) {
	const op = "services.user.ByEmail"

	user, err := s.userRepo.ByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil
}

func (s *Service) IsAdmin(ctx context.Context, userId string) (bool, error) {
	const op = "services.user.IsAdmin"

	uuid, err := uuid.Parse(userId)
	if err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}

	member, err := s.userRepo.StaffMemberByUuid(ctx, uuid)
	if err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}

	for _, permission := range member.Permissions {
		if permission == entity.Admin {
			return true, nil
		}
	}

	return false, nil
}

func (s *Service) Users(ctx context.Context) ([]entity.User, error) {
	return s.userRepo.Users(ctx)
}
