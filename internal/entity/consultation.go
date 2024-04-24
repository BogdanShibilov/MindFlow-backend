package entity

import (
	"time"

	"github.com/google/uuid"
)

type Consultation struct {
	Uuid                    uuid.UUID `db:"uuid"`
	ExpertUuid              uuid.UUID `db:"expert_uuid"`
	MenteeUuid              uuid.UUID `db:"mentee_uuid"`
	ConsultationApplication `db:"-"`
}

type ConsultationApplication struct {
	ConsultationUuid uuid.UUID `db:"consultation_uuid"`
	Status           Status    `db:"status"`
	MenteeQuestions  string    `db:"mentee_questions"`
	SubmittedAt      time.Time `db:"submitted_at"`
}

type ConsultationMeeting struct {
	Uuid             uuid.UUID `db:"uuid"`
	ConsultationUuid uuid.UUID `db:"consultation_uuid"`
	StartTime        time.Time `db:"start_time"`
	Link             string    `db:"link"`
}
