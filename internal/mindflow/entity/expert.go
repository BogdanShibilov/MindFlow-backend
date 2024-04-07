package entity

import (
	"time"

	"github.com/google/uuid"
)

type TimeUnit string

type ExpertInfo struct {
	UserUuid               uuid.UUID
	Position               string
	ChargePerHour          int
	ExperienceDescription  string
	ExpertiseAtDescription string
}

type ExpertApplication struct {
	ExpertInfo  *ExpertInfo
	SubmittedAt time.Time
	ApprovedAt  time.Time
	IsApproved  bool
}
