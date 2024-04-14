package entity

import (
	"fmt"

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
	Surname               string
	PhoneNumber           string
	ProfessionalField     string
	ExperienceDescription string
}

func (u *UserDetails) FullName() string {
	return fmt.Sprintf("%s %s", u.Name, u.Surname)
}
