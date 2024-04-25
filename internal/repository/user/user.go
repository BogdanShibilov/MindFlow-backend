package userrepo

import (
	"context"
	"errors"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/bogdanshibilov/mindflowbackend/internal/db/postgres"
	"github.com/bogdanshibilov/mindflowbackend/internal/entity"
)

type Repo struct {
	Db postgres.Db
}

func (r *Repo) CreateUser(ctx context.Context, user *entity.User) error {
	const op = "repository.user.CreateUser"

	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	insertUserSql, insertUserArgs, err := psql.Insert("users").
		Columns("username", "pass_hash").
		Values(user.Username, user.PassHash).
		ToSql()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	tx, err := r.Db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback(ctx)
		} else {
			_ = tx.Commit(ctx)
		}
	}()

	_, err = tx.Exec(ctx, insertUserSql, insertUserArgs...)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	selectInsertedUserUuidSql, selectInsertedUserUuidArgs, err := psql.Select("uuid").
		From("users").
		Where("username in (?)", user.Username).
		ToSql()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	var uuid string
	err = tx.QueryRow(ctx, selectInsertedUserUuidSql, selectInsertedUserUuidArgs...).Scan(&uuid)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	insertProfileSql, insertProfileArgs, err := psql.Insert("user_profiles").
		Columns(
			"user_uuid",
			"name",
			"email",
			"phone",
			"professional_field",
			"experience_description",
		).
		Values(
			uuid,
			user.Name,
			user.Email,
			user.Phone,
			user.ProfessionalField,
			user.ExperienceDescription,
		).
		ToSql()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	_, err = tx.Exec(ctx, insertProfileSql, insertProfileArgs...)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (r *Repo) UpdateCredentials(
	ctx context.Context,
	creds *entity.UserCredentials,
	uuid uuid.UUID,
) error {
	const op = "repository.user.UpdateCredentials"

	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	sql, args, err := psql.Update("users").
		SetMap(
			sq.Eq{
				"username":  creds.Username,
				"pass_hash": creds.Username,
			},
		).
		Where("uuid IN (?)", uuid).
		ToSql()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	_, err = r.Db.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (r *Repo) UpdateProfile(
	ctx context.Context,
	profile *entity.UserProfile,
	uuid uuid.UUID,
) error {
	const op = "repository.user.UpdateProfile"

	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	sql, args, err := psql.Update("user_profiles").
		SetMap(
			sq.Eq{
				"name":                   profile.Name,
				"email":                  profile.Email,
				"phone":                  profile.Phone,
				"professional_field":     profile.ProfessionalField,
				"experience_description": profile.ExperienceDescription,
			},
		).
		Where("user_uuid IN (?)", uuid).
		ToSql()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	_, err = r.Db.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (r *Repo) ByUuid(ctx context.Context, uuid uuid.UUID) (*entity.User, error) {
	const op = "repository.user.ByUuid"

	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	sql, args, err := psql.Select(
		"uuid",
		"username",
		"pass_hash",
		"roles",
		"name",
		"email",
		"phone",
		"professional_field",
		"experience_description",
	).
		From("users").
		InnerJoin("user_profiles ON users.uuid = user_profiles.user_uuid").
		Where("user.uuid IN (?)", uuid).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	rows, err := r.Db.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	user, err := pgx.CollectOneRow(rows, pgx.RowToStructByNameLax[entity.User])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("%s: %w", op, ErrUserNotFound)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &user, nil
}

func (r *Repo) ByEmail(ctx context.Context, email string) (*entity.User, error) {
	const op = "repository.user.ByEmail"

	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	sql, args, err := psql.Select(
		"uuid",
		"username",
		"pass_hash",
		"roles",
		"name",
		"email",
		"phone",
		"professional_field",
		"experience_description",
	).
		From("users").InnerJoin("user_profiles ON users.uuid = user_profiles.user_uuid").
		Where("email in (?)", email).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	rows, err := r.Db.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	user, err := pgx.CollectOneRow(rows, pgx.RowToStructByNameLax[entity.User])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("%s: %w", op, ErrUserNotFound)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &user, nil
}

func (r *Repo) StaffMemberByUuid(ctx context.Context, uuid uuid.UUID) (*entity.StaffMember, error) {
	const op = "repository.user.StaffMemberByUuid"

	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	sql, args, err := psql.Select("user_uuid", "permissions").
		From("staff").
		Where("user_uuid IN (?)", uuid).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	rows, err := r.Db.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	member, err := pgx.CollectOneRow(rows, pgx.RowToStructByNameLax[entity.StaffMember])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("%s: %w", op, ErrUserNotFound)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &member, nil
}
