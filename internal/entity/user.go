package entity

import "github.com/google/uuid"

type User struct {
	Uuid            uuid.UUID `db:"uuid"`
	Roles           []string  `db:"roles"`
	UserCredentials `db:"-"`
	UserProfile     `db:"-"`
}

type UserCredentials struct {
	Username string `db:"username"`
	PassHash []byte `db:"pass_hash"`
}

type UserProfile struct {
	Email                 string `db:"email"`
	Phone                 string `db:"phone"`
	ProfessionalField     string `db:"professional_field"`
	ExperienceDescription string `db:"experience_description"`
}

type StaffMember struct {
	UserUuid    uuid.UUID    `db:"user_uuid"`
	Permissions []Permission `db:"permissions"`
}
