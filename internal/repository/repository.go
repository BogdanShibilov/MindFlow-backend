package repository

import (
	"github.com/bogdanshibilov/mindflowbackend/internal/db/postgres"
	consultationrepo "github.com/bogdanshibilov/mindflowbackend/internal/repository/consultation"
	expertrepo "github.com/bogdanshibilov/mindflowbackend/internal/repository/expert"
	userrepo "github.com/bogdanshibilov/mindflowbackend/internal/repository/user"
)

func NewUser(db *postgres.Db) *userrepo.Repo {
	return &userrepo.Repo{
		Db: *db,
	}
}

func NewExpert(db *postgres.Db) *expertrepo.Repo {
	return &expertrepo.Repo{
		Db: *db,
	}
}

func NewConsultation(db *postgres.Db) *consultationrepo.Repo {
	return &consultationrepo.Repo{
		Db: *db,
	}
}
