package expert

import (
	"context"

	"github.com/google/uuid"

	"github.com/bogdanshibilov/mindflowbackend/internal/mindflow/entity"
)

type service struct {
	expertSaver         ExpertSaver
	expertProvider      ExpertProvider
	applicationSaver    ApplicationSaver
	applicationProvider ApplicationProvider
}

type ExpertSaver interface {
	SaveExpertInfo(ctx context.Context, expertInfo *entity.ExpertInfo) error
}

type ExpertProvider interface {
	ExpertByUuid(ctx context.Context, uuid *uuid.UUID) (*entity.ExpertInfo, error)
}

type ApplicationSaver interface {
	SaveApplication(ctx context.Context, application *entity.ExpertApplication) error
}

type ApplicationProvider interface {
	ApplicationByUuid(ctx context.Context, uuid *uuid.UUID) (*entity.ExpertApplication, error)
}

func New(
	expertSaver ExpertSaver,
	expertProvider ExpertProvider,
	applicationSaver ApplicationSaver,
	applicationProvider ApplicationProvider,
) *service {
	return &service{
		expertSaver:         expertSaver,
		expertProvider:      expertProvider,
		applicationSaver:    applicationSaver,
		applicationProvider: applicationProvider,
	}
}
