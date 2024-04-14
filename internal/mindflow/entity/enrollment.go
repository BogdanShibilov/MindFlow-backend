package entity

import (
	"time"

	"github.com/google/uuid"
)

type Enrollment struct {
	Uuid       uuid.UUID
	MentorUuid uuid.UUID
	MenteeUuid uuid.UUID
}

type Meeting struct {
	Uuid           uuid.UUID
	EnrollmentUuid uuid.UUID
	Link           string
	StartTime      time.Time
	EndTime        time.Time
}
