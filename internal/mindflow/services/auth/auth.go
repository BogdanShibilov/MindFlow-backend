package auth

import (
	"context"
	"errors"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/bogdanshibilov/mindflowbackend/internal/mindflow/entity"
	"github.com/bogdanshibilov/mindflowbackend/internal/mindflow/lib/jwt"
	"github.com/google/uuid"
)

type Auth struct {
	userSaver    UserSaver
	userProvider UserProvider
	tokenTTL     time.Duration
	secret       string
}

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
)

type UserSaver interface {
	SaveUser(ctx context.Context, user *entity.User) error
	SaveUserDetails(ctx context.Context, userDetails *entity.UserDetails) error
	UpdateUserDetails(ctx context.Context, userDetails *entity.UserDetails) error
}

type UserProvider interface {
	UserByEmail(ctx context.Context, email string) (*entity.User, error)
	UserDetailsByUserUuid(ctx context.Context, uuid uuid.UUID) (*entity.UserDetails, error)
}

func New(
	userSaver UserSaver,
	userProvider UserProvider,
	tokenTTL time.Duration,
	secret string,
) *Auth {
	return &Auth{
		userSaver:    userSaver,
		userProvider: userProvider,
		tokenTTL:     tokenTTL,
		secret:       secret,
	}
}

func (a *Auth) RegisterNewUser(ctx context.Context, email, pass string, name string) error {
	const op = "Auth.RegisterNewUser"

	passHash, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	newUser := &entity.User{
		Email:    email,
		PassHash: passHash,
		Name:     name,
	}

	err = a.userSaver.SaveUser(ctx, newUser)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	newUser, err = a.userProvider.UserByEmail(ctx, newUser.Email)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	newUserDetails := &entity.UserDetails{
		UserUuid:              newUser.Uuid,
		PhoneNumber:           "",
		ProfessionalField:     "",
		ExperienceDescription: "",
	}
	err = a.userSaver.SaveUserDetails(ctx, newUserDetails)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (a *Auth) Login(ctx context.Context, email, pass string) (string, error) {
	const op = "Auth.Login"

	user, err := a.userProvider.UserByEmail(ctx, email)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	if err := bcrypt.CompareHashAndPassword(user.PassHash, []byte(pass)); err != nil {
		return "", fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
	}

	token, err := jwt.NewAccessToken(user, a.secret, a.tokenTTL)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return token, nil
}

func (a *Auth) UpdateUserDetails(
	ctx context.Context,
	id string,
	phoneNumber string,
	professionalField string,
	experienceDescription string,
) error {
	const op = "Auth.UpdateUserDetails"

	uuid, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	newUserDetails := &entity.UserDetails{
		UserUuid:              uuid,
		PhoneNumber:           phoneNumber,
		ProfessionalField:     professionalField,
		ExperienceDescription: experienceDescription,
	}

	err = a.userSaver.UpdateUserDetails(ctx, newUserDetails)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (a *Auth) UserDetailsByUserUuid(ctx context.Context, id string) (*entity.UserDetails, error) {
	const op = "Auth.UserDetailsByUserUuid"

	uuid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	userDetails, err := a.userProvider.UserDetailsByUserUuid(ctx, uuid)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return userDetails, nil
}
