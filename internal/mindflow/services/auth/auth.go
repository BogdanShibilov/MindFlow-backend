package auth

import (
	"context"
	"log/slog"
	"time"

	"github.com/bogdanshibilov/mindflowbackend/internal/mindflow/entity"
)

type Auth struct {
	log          *slog.Logger
	userSaver    UserSaver
	UserProvider UserProvider
	tokenTTL     time.Duration
}

type UserSaver interface {
	SaveUser(
		ctx context.Context,
		email string,
		passHash []byte,
	) (guid string, err error)
}

type UserProvider interface {
	User(ctx context.Context, email string) (entity.User, error)
}

func New(
	log *slog.Logger,
	userSaver UserSaver,
	userProvider UserProvider,
	tokenTTL time.Duration,
) *Auth {
	return &Auth{
		log:          log,
		userSaver:    userSaver,
		UserProvider: userProvider,
		tokenTTL:     tokenTTL,
	}
}
