package auth

import (
	"context"
	"errors"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/bogdanshibilov/mindflowbackend/internal/mindflow/entity"
	"github.com/bogdanshibilov/mindflowbackend/internal/mindflow/lib/jwt"
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
}

type UserProvider interface {
	UserByEmail(ctx context.Context, email string) (*entity.User, error)
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

func (a *Auth) RegisterNewUser(ctx context.Context, email, pass string) error {
	const op = "Auth.RegisterNewUser"

	passHash, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	newUser := &entity.User{
		Email:    email,
		PassHash: passHash,
	}

	err = a.userSaver.SaveUser(ctx, newUser)
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
