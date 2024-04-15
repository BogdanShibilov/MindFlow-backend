package postgres

import (
	"context"
	"errors"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/bogdanshibilov/mindflowbackend/internal/mindflow/entity"
	"github.com/bogdanshibilov/mindflowbackend/internal/mindflow/storage"
)

func (s *Storage) SaveUser(ctx context.Context, user *entity.User) error {
	const op = "storage.postgres.SaveUser"
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	sql, args, err := psql.Insert("users").
		Columns("email", "pass_hash", "name").
		Values(user.Email, user.PassHash, user.Name).
		ToSql()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	_, err = s.conn.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) UserByEmail(ctx context.Context, email string) (*entity.User, error) {
	const op = "storage.postgres.UserByEmail"
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	sql, args, err := psql.Select("uuid", "email", "pass_hash", "name").
		From("users").
		Where("email IN (?)", email).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	var user entity.User
	row := s.conn.QueryRow(ctx, sql, args...)
	err = row.Scan(&user.Uuid, &user.Email, &user.PassHash, &user.Name)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("%s: %w", op, storage.ErrEntityNotFound)
		}

		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &user, nil
}

func (s *Storage) SaveUserDetails(ctx context.Context, userDetails *entity.UserDetails) error {
	const op = "storage.postgres.SaveUserDetails"
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	sql, args, err := psql.Insert("user_details").
		Columns("user_uuid", "phone_number",
			"professional_field", "experience_description").
		Values(userDetails.UserUuid, userDetails.PhoneNumber,
			userDetails.ProfessionalField, userDetails.ExperienceDescription).
		ToSql()

	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	_, err = s.conn.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) UpdateUserDetails(ctx context.Context, userDetails *entity.UserDetails) error {
	const op = "storage.postgres.UpdateUserDetails"
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	sql, args, err := psql.Update("user_details").
		SetMap(
			sq.Eq{
				"phone_number":           userDetails.PhoneNumber,
				"professional_field":     userDetails.ProfessionalField,
				"experience_description": userDetails.ProfessionalField,
			},
		).
		Where("user_uuid IN (?)", userDetails.UserUuid).
		ToSql()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	_, err = s.conn.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) UserDetailsByUserUuid(ctx context.Context, uuid uuid.UUID) (*entity.UserDetails, error) {
	const op = "storage.postgres.UserDetailsByUserUuid"
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	sql, args, err := psql.Select("user_uuid", "phone_number",
		"professional_field", "experience_description").
		From("user_details").
		Where("user_uuid IN (?)", uuid).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	var userDetails entity.UserDetails
	row := s.conn.QueryRow(ctx, sql, args...)
	err = row.Scan(&userDetails.UserUuid,
		&userDetails.PhoneNumber, &userDetails.ProfessionalField,
		&userDetails.ExperienceDescription)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("%s: %w", op, storage.ErrEntityNotFound)
		}

		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &userDetails, nil
}
