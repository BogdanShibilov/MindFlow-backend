package authservice

import (
	"context"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"

	jwtservice "github.com/bogdanshibilov/mindflowbackend/internal/services/jwt"
	userservice "github.com/bogdanshibilov/mindflowbackend/internal/services/user"
)

type Service struct {
	users    *userservice.Service
	secret   string
	tokenTTL time.Duration
}

func New(users *userservice.Service, secret string, tokenTTL time.Duration) *Service {
	return &Service{
		users:    users,
		secret:   secret,
		tokenTTL: tokenTTL,
	}
}

func (s *Service) Register(
	ctx context.Context,
	username string,
	password string,
	email string,
	phone string,
	professionalField string,
	experienceDescription string,
) error {
	const op = "services.auth.Register"

	passHash, err := bcrypt.
		GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	err = s.users.CreateUser(
		ctx,
		username,
		string(passHash),
		email,
		phone,
		professionalField,
		experienceDescription,
	)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Service) LoginByEmail(
	ctx context.Context,
	email string,
	password string,
) (token string, err error) {
	const op = "services.auth.Login"

	user, err := s.users.ByEmail(ctx, email)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	if err := bcrypt.CompareHashAndPassword(user.PassHash, []byte(password)); err != nil {
		return "", fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
	}

	accessToken, err := jwtservice.NewAccessToken(
		user.Uuid.String(),
		user.Email,
		user.Roles,
		s.secret,
		s.tokenTTL,
	)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return accessToken, nil
}
