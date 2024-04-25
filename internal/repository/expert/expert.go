package expertrepo

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

func (r *Repo) CreateExpert(ctx context.Context, expert *entity.Expert) error {
	const op = "repository.expert.CreateExpert"

	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	insertInfoSql, insertInfoArgs, err := psql.Insert("expert_information").
		Columns(
			"user_uuid",
			"price",
			"help_description",
		).
		Values(
			expert.UserUuid,
			expert.Price,
			expert.HelpDescription,
		).
		ToSql()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	insertApplicationSql, insertApplicationArgs, err := psql.Insert("expert_application").
		Columns(
			"user_uuid",
			"status",
		).
		Values(
			expert.UserUuid,
			expert.Status,
		).
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

	_, err = tx.Exec(ctx, insertInfoSql, insertInfoArgs...)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	_, err = tx.Exec(ctx, insertApplicationSql, insertApplicationArgs...)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (r *Repo) ByUuid(ctx context.Context, uuid uuid.UUID) (*entity.Expert, error) {
	const op = "repository.expert.ByUuid"

	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	sql, args, err := psql.Select(
		"expert_information.user_uuid AS user_uuid",
		"price",
		"help_description",
		"status",
		"submitted_at",
		"email",
		"name",
		"phone",
		"professional_field",
		"experience_description",
	).
		From("expert_information").
		InnerJoin("expert_application ON expert_information.user_uuid = expert_application.user_uuid").
		InnerJoin("user_profiles ON expert_application.user_uuid = user_profiles.user_uuid").
		Where("expert_information.user_uuid IN (?)", uuid).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	rows, err := r.Db.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	expert, err := pgx.CollectOneRow(rows, pgx.RowToStructByNameLax[entity.Expert])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("%s: %w", op, ErrExpertNotFound)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &expert, nil
}

// Option statuses: Pending, Accepted, Rejected, by default: All
func (r *Repo) Experts(ctx context.Context, opts ...SelectExpertsOption) ([]entity.Expert, error) {
	const op = "repository.expert.Experts"

	options := newDefaultExpertsSelectOptions()
	for _, opt := range opts {
		opt(options)
	}

	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	expertsQuery := psql.Select(
		"expert_information.user_uuid AS user_uuid",
		"price",
		"help_description",
		"status",
		"submitted_at",
		"email",
		"name",
		"phone",
		"professional_field",
		"experience_description",
	).
		From("expert_information").
		InnerJoin("expert_application ON expert_information.user_uuid = expert_application.user_uuid").
		InnerJoin("user_profiles ON expert_application.user_uuid = user_profiles.user_uuid")
	if options.status != AllStatus {
		expertsQuery = expertsQuery.Where("status IN (?)", options.status)
	}
	sql, args, err := expertsQuery.ToSql()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	rows, err := r.Db.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	experts, err := pgx.CollectRows(rows, pgx.RowToStructByNameLax[entity.Expert])
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return experts, err
}

func (r *Repo) UpdateExpertApplication(
	ctx context.Context,
	application *entity.ExpertApplication,
	userUuid uuid.UUID,
) error {
	const op = "repository.expert.UpdateExpertApplication"

	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	sql, args, err := psql.Update("expert_application").
		SetMap(sq.Eq{
			"status": application.Status,
		},
		).
		Where("user_uuid IN (?)", userUuid).
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
