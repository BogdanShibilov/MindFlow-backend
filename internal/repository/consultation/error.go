package consultationrepo

import "errors"

var (
	ErrConsultationNotFound = errors.New("consultation not found")
	ErrMeetingFKViolation   = errors.New("tried to create a meeting for non existing consultation")
)
