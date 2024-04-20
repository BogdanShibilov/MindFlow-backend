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

func (s *Storage) SaveExpertInfo(ctx context.Context, info *entity.ExpertInfo) error {
	const op = "storage.postgres.SaveExpertInfo"
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	sql, args, err := psql.Insert("expert_info").
		Columns("user_uuid", "charge_per_hour",
			"expertise_at_description").
		Values(info.UserUuid, info.ChargePerHour,
			info.ExpertiseAtDescription).
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

func (s *Storage) UpdateExpertInfo(ctx context.Context, info *entity.ExpertInfo) error {
	const op = "storage.postgres.UpdateExpertInfo"
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	sql, args, err := psql.Update("expert_info").
		SetMap(
			sq.Eq{
				"charge_per_hour":          info.ChargePerHour,
				"expertise_at_description": info.ExpertiseAtDescription,
				"is_approved":              info.IsApproved,
			},
		).
		Where("user_uuid IN (?)", info.UserUuid).
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

func (s *Storage) ExpertByUuid(ctx context.Context, uuid *uuid.UUID) (*entity.ExpertInfo, error) {
	const op = "storage.postgres.ExpertByUuid"
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	sql, args, err := psql.Select("user_uuid", "charge_per_hour",
		"expertise_at_description", "submitted_at",
		"is_approved").
		From("expert_info").
		Where("user_uuid IN (?)", uuid).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	var expertInfo entity.ExpertInfo
	row := s.conn.QueryRow(ctx, sql, args...)
	err = row.Scan(&expertInfo.UserUuid, &expertInfo.ChargePerHour,
		&expertInfo.ExpertiseAtDescription,
		&expertInfo.SubmittedAt, &expertInfo.IsApproved)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("%s: %w", op, storage.ErrEntityNotFound)
		}

		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &expertInfo, nil
}

func (s *Storage) ExpertInfo(ctx context.Context) ([]entity.ExpertInfo, error) {
	const op = "storage.postgres.ExpertInfo"
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	sql, args, err := psql.Select("user_uuid", "charge_per_hour",
		"expertise_at_description", "submitted_at",
		"is_approved").
		From("expert_info").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	rows, err := s.conn.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	expertInfo, err := pgx.CollectRows(rows, pgx.RowToStructByPos[entity.ExpertInfo])
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return expertInfo, nil
}

func (s *Storage) ApprovedExpertInfo(ctx context.Context) ([]entity.ExpertInfo, error) {
	const op = "storage.postgres.ApprovedExpertInfo"
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	sql, args, err := psql.Select("user_uuid", "charge_per_hour",
		"expertise_at_description", "submitted_at",
		"is_approved").
		From("expert_info").
		Where("is_approved = TRUE").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	rows, err := s.conn.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	expertInfo, err := pgx.CollectRows(rows, pgx.RowToStructByPos[entity.ExpertInfo])
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return expertInfo, nil
}

func (s *Storage) NonApprovedExpertInfo(ctx context.Context) ([]entity.ExpertInfo, error) {
	const op = "storage.postgres.NonApprovedExpertInfo"
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	sql, args, err := psql.Select("user_uuid", "charge_per_hour",
		"expertise_at_description", "submitted_at",
		"is_approved").
		From("expert_info").
		Where("is_approved = FALSE").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	rows, err := s.conn.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	expertInfo, err := pgx.CollectRows(rows, pgx.RowToStructByPos[entity.ExpertInfo])
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return expertInfo, nil
}
