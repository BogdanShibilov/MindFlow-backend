package entity

import "github.com/google/uuid"

type User struct {
	Uuid     uuid.UUID
	Email    string
	PassHash []byte
}
