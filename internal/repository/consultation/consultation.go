package consultationrepo

import (
	"context"
	"errors"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	"github.com/bogdanshibilov/mindflowbackend/internal/db/postgres"
	"github.com/bogdanshibilov/mindflowbackend/internal/entity"
)

type Repo struct {
	Db postgres.Db
}

func (r *Repo) CreateConsultation(ctx context.Context, consult *entity.Consultation) error {
	const op = "repository.consultation.CreateConsultation"

	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	insertConsultSql, insertConsultArgs, err := psql.Insert("consultation").
		Columns(
			"expert_uuid",
			"mentee_uuid",
		).
		Values(
			consult.ExpertUuid,
			consult.MenteeUuid,
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

	_, err = tx.Exec(ctx, insertConsultSql, insertConsultArgs...)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	selectInsertedConsultSql, selectInsertedConsultArgs, err := psql.Select("uuid").
		From("consultation").
		Where(
			"expert_uuid IN (?) AND mentee_uuid IN (?)",
			consult.ExpertUuid,
			consult.MenteeUuid,
		).
		ToSql()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	var consultUuid uuid.UUID
	err = tx.QueryRow(ctx, selectInsertedConsultSql, selectInsertedConsultArgs...).Scan(&consultUuid)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	insertApplicationSql, insertApplicationArgs, err := psql.Insert("consultation_application").
		Columns(
			"consultation_uuid",
			"mentee_questions",
		).
		Values(
			consultUuid,
			consult.MenteeQuestions,
		).
		ToSql()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	_, err = tx.Exec(ctx, insertApplicationSql, insertApplicationArgs...)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (r *Repo) Consultations(ctx context.Context) ([]entity.Consultation, error) {
	const op = "repository.consultation.Consultations"

	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	sql, args, err := psql.Select(
		"consultation.uuid AS uuid",
		"mentee_uuid",
		"status",
		"mentee_questions",
		"submitted_at",
	).
		From("consultation").
		InnerJoin("consultation_application ON consultation.uuid = consultation_application.consultation_uuid").
		Where("status IN (?)", entity.Pending).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	rows, err := r.Db.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	applications, err := pgx.CollectRows(rows, pgx.RowToStructByNameLax[entity.Consultation])
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return applications, err
}

func (r *Repo) ByUuid(ctx context.Context, uuid uuid.UUID) (*entity.Consultation, error) {
	const op = "repository.consultation.ByUuid"

	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	sql, args, err := psql.Select(
		"consultation.uuid AS uuid",
		"expert_uuid",
		"mentee_uuid",
		"status",
		"mentee_questions",
		"submitted_at",
	).
		From("consultation").
		InnerJoin("consultation_application ON consultation.uuid = consultation_application.consultation_uuid").
		Where("uuid IN (?)", uuid).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	rows, err := r.Db.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	consultation, err := pgx.CollectOneRow(rows, pgx.RowToStructByNameLax[entity.Consultation])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("%s: %w", op, ErrConsultationNotFound)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &consultation, nil
}

func (r *Repo) ByPersonUuid(ctx context.Context, uuid uuid.UUID, opts ...ByPersonUuidOption) ([]entity.Consultation, error) {
	const op = "repository.consultation.ByPersonUuid"

	options := newDefaultSelectByPersonUuidOptions()
	for _, opt := range opts {
		opt(options)
	}

	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	sql, args, err := psql.Select(
		"consultation.uuid AS uuid",
		"expert_uuid",
		"mentee_uuid",
		"status",
		"mentee_questions",
		"submitted_at",
	).
		From("consultation").
		InnerJoin("consultation_application ON consultation.uuid = consultation_application.consultation_uuid").
		Where(string(options.whoseUuid)+" IN (?)", uuid).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	rows, err := r.Db.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	consultations, err := pgx.CollectRows(rows, pgx.RowToStructByNameLax[entity.Consultation])
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return consultations, nil
}

func (r *Repo) MeetingsByConsultationUuid(ctx context.Context, uuid uuid.UUID) ([]entity.ConsultationMeeting, error) {
	const op = "repository.consultation.MeetingsByConsultationUuid"

	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	sql, args, err := psql.Select(
		"uuid",
		"consultation_uuid",
		"start_time",
		"link",
	).
		From("consultation_meeting").
		Where("consultation_uuid IN (?)", uuid).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	rows, err := r.Db.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	meetings, err := pgx.CollectRows(rows, pgx.RowToStructByNameLax[entity.ConsultationMeeting])
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return meetings, nil
}

func (r *Repo) CreateMeeting(ctx context.Context, meeting *entity.ConsultationMeeting) error {
	const op = "repository.consultation.CreateMeeting"

	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	sql, args, err := psql.Insert("consultation_meeting").
		Columns(
			"consultation_uuid",
			"start_time",
			"link",
		).
		Values(
			meeting.ConsultationUuid,
			meeting.StartTime,
			meeting.Link,
		).
		ToSql()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	_, err = r.Db.Exec(ctx, sql, args...)
	if err != nil {
		var pgError *pgconn.PgError
		if errors.As(err, &pgError) {
			if pgError.Code == pgerrcode.ForeignKeyViolation {
				return fmt.Errorf("%s: %w", op, ErrMeetingFKViolation)
			}
		}
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (r *Repo) UpdateApplication(ctx context.Context, application *entity.ConsultationApplication) error {
	const op = "repository.consultation.UpdateApplication"

	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	sql, args, err := psql.Update("consultation_application").
		SetMap(
			sq.Eq{
				"status":           application.Status,
				"mentee_questions": application.MenteeQuestions,
			},
		).
		Where("consultation_uuid IN (?)", application.ConsultationUuid).
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

func (r *Repo) UpdateApplicationStatus(ctx context.Context, uuid uuid.UUID, status entity.Status) error {
	const op = "repository.consultation.UpdateApplicationStatus"

	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	sql, args, err := psql.Update("consultation_application").
		SetMap(
			sq.Eq{
				"status": status,
			},
		).
		Where("consultation_uuid IN (?)", uuid).
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

func (r *Repo) DoesExist(ctx context.Context, menteeUuid, expertUuid uuid.UUID) (bool, error) {
	const op = "repository.consultation.DoesExist"

	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	sql, args, err := psql.Select(
		"Count(*)",
	).
		From("consultation").
		Where("expert_uuid IN (?)", expertUuid).
		Where("mentee_uuid IN (?)", menteeUuid).
		ToSql()
	if err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}

	var count int
	err = r.Db.QueryRow(ctx, sql, args...).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}

	if count > 0 {
		return true, nil
	} else {
		return false, nil
	}
}
