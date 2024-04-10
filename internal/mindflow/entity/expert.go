package entity

import (
	"time"

	"github.com/google/uuid"
)

type ExpertInfo struct {
	UserUuid               uuid.UUID
	Position               string
	ChargePerHour          int
	ExperienceDescription  string
	ExpertiseAtDescription string
	SubmittedAt            time.Time
	IsApproved             bool
}
