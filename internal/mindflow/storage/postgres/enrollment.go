package postgres

import (
	"context"
	"errors"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/bogdanshibilov/mindflowbackend/internal/mindflow/entity"
	"github.com/bogdanshibilov/mindflowbackend/internal/mindflow/storage"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func (s *Storage) SaveEnrollment(ctx context.Context, enrollment *entity.Enrollment) error {
	const op = "storage.postgres.SaveEnrollment"
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	sql, args, err := psql.Insert("enrollments").
		Columns("mentor_uuid", "mentee_uuid",
			"is_approved", "mentee_questions").
		Values(enrollment.MentorUuid, enrollment.MenteeUuid,
			enrollment.IsApproved, enrollment.MenteeQuestions).
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

func (s *Storage) UpdateEnrollment(ctx context.Context, enrollment *entity.Enrollment) error {
	const op = "storage.postgres.UpdateEnrollment"
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	sql, args, err := psql.Update("enrollments").
		SetMap(
			sq.Eq{
				"mentor_uuid":      enrollment.MentorUuid,
				"mentee_uuid":      enrollment.MenteeUuid,
				"is_approved":      enrollment.IsApproved,
				"mentee_questions": enrollment.MenteeQuestions,
			},
		).
		Where("uuid IN (?)", enrollment.Uuid).
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

func (s *Storage) EnrollmentByUuid(ctx context.Context, uuid uuid.UUID) (*entity.Enrollment, error) {
	const op = "storage.postgres.EnrollmentByUuid"
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	sql, args, err := psql.Select("uuid", "mentor_uuid", "mentee_uuid",
		"is_approved", "mentee_questions").
		From("enrollments").
		Where("uuid IN (?)", uuid).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	var enrollment entity.Enrollment
	row := s.conn.QueryRow(ctx, sql, args...)
	err = row.Scan(&enrollment.Uuid, &enrollment.MentorUuid, &enrollment.MenteeUuid,
		&enrollment.IsApproved, &enrollment.MenteeQuestions)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("%s: %w", op, storage.ErrEntityNotFound)
		}

		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &enrollment, nil
}

func (s *Storage) EnrollmentsByMenteeUuid(ctx context.Context, uuid uuid.UUID) ([]entity.Enrollment, error) {
	const op = "storage.postgres.EnrollmentsByMenteeUuid"
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	sql, args, err := psql.Select("uuid", "mentor_uuid", "mentee_uuid",
		"is_approved", "mentee_questions").
		From("enrollments").
		Where("mentee_uuid = (?)", uuid).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	rows, err := s.conn.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	enrollments, err := pgx.CollectRows(rows, pgx.RowToStructByPos[entity.Enrollment])
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return enrollments, nil
}

func (s *Storage) EnrollmentsByMentorUuid(ctx context.Context, uuid uuid.UUID) ([]entity.Enrollment, error) {
	const op = "storage.postgres.EnrollmentsByMentorUuid"
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	sql, args, err := psql.Select("uuid", "mentor_uuid", "mentee_uuid",
		"is_approved", "mentee_questions").
		From("enrollments").
		Where("mentor_uuid = (?)", uuid).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	rows, err := s.conn.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	enrollments, err := pgx.CollectRows(rows, pgx.RowToStructByPos[entity.Enrollment])
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return enrollments, nil
}
