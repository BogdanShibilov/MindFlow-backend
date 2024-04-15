package entity

import (
	"github.com/google/uuid"
)

type User struct {
	Uuid     uuid.UUID
	Email    string
	PassHash []byte
	Roles    []string
}

type UserDetails struct {
	UserUuid              uuid.UUID
	Name                  string
	PhoneNumber           string
	ProfessionalField     string
	ExperienceDescription string
}
