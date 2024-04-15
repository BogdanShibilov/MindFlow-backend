package enrollment

import (
	"context"
	"fmt"

	"github.com/bogdanshibilov/mindflowbackend/internal/mindflow/entity"
	"github.com/google/uuid"
)

type ByWho string

const (
	Mentor ByWho = "mentor"
	Mentee ByWho = "mentee"
)

type Service struct {
	enrollmentSaver    EnrollmentSaver
	enrollmentProvider EnrollmentProvider
}

type EnrollmentSaver interface {
	SaveEnrollment(ctx context.Context, enrollment *entity.Enrollment) error
	UpdateEnrollment(ctx context.Context, enrollment *entity.Enrollment) error
}

type EnrollmentProvider interface {
	EnrollmentByUuid(ctx context.Context, uuid uuid.UUID) (*entity.Enrollment, error)
	EnrollmentsByMenteeUuid(ctx context.Context, uuid uuid.UUID) ([]entity.Enrollment, error)
	EnrollmentsByMentorUuid(ctx context.Context, uuid uuid.UUID) ([]entity.Enrollment, error)
}

func New(
	enrollmentSaver EnrollmentSaver,
	enrollmentProvider EnrollmentProvider,
) *Service {
	return &Service{
		enrollmentSaver:    enrollmentSaver,
		enrollmentProvider: enrollmentProvider,
	}
}

func (s *Service) CreateEnrollment(
	ctx context.Context,
	mentorId string,
	menteeId string,
	menteeQuestions string,
) error {
	const op = "Expert.CreateEnrollment"

	mentorUuid, err := uuid.Parse(mentorId)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	menteeUuid, err := uuid.Parse(menteeId)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	enrollment := &entity.Enrollment{
		MentorUuid:      mentorUuid,
		MenteeUuid:      menteeUuid,
		IsApproved:      false,
		MenteeQuestions: menteeQuestions,
	}

	err = s.enrollmentSaver.SaveEnrollment(ctx, enrollment)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

// Gets enrollments by someone's id
// EnrollmentsByMemberId(ctx, id, enrollment.Mentor) will look all enrollments by mentor's id
// EnrollmentsByMemberId(ctx, id, enrollment.Mentee) will look all enrollments by mentee's id
func (s *Service) EnrollmentsByMemberId(
	ctx context.Context,
	id string,
	byWho ByWho,
) ([]entity.Enrollment, error) {
	const op = "Expert.EnrollmentsByMentorId"

	uuid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	var enrollments []entity.Enrollment
	if byWho == Mentor {
		enrollments, err = s.enrollmentProvider.EnrollmentsByMentorUuid(ctx, uuid)
	} else {
		enrollments, err = s.enrollmentProvider.EnrollmentsByMenteeUuid(ctx, uuid)
	}
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return enrollments, nil
}

func (s *Service) ApproveEnrollmentById(ctx context.Context, id string) error {
	const op = "Expert.ApproveEnrollmentById"

	uuid, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	enrollment, err := s.enrollmentProvider.EnrollmentByUuid(ctx, uuid)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	enrollment.IsApproved = true
	err = s.enrollmentSaver.UpdateEnrollment(ctx, enrollment)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
