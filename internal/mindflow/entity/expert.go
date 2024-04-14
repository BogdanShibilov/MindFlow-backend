package entity

import (
	"time"

	"github.com/google/uuid"
)

type ExpertInfo struct {
	UserUuid               uuid.UUID
	ChargePerHour          int
	ExpertiseAtDescription string
	SubmittedAt            time.Time
	IsApproved             bool
}
