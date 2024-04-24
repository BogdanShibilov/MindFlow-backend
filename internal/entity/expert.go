package entity

import (
	"time"

	"github.com/google/uuid"
)

type Expert struct {
	UserUuid          uuid.UUID `db:"user_uuid"`
	UserProfile       `db:"-"`
	ExpertInformation `db:"-"`
	ExpertApplication `db:"-"`
}

type ExpertInformation struct {
	Price           int    `db:"price"`
	HelpDescription string `db:"help_description"`
}

type ExpertApplication struct {
	Status      Status    `db:"status"`
	SubmittedAt time.Time `db:"submitted_at"`
}
