package expertrepo

import "github.com/bogdanshibilov/mindflowbackend/internal/entity"

const (
	AllStatus entity.Status = -1
)

const (
	_defaultExpertStatus = AllStatus
)

type expertSelectOptions struct {
	status entity.Status
}

func newDefaultExpertsSelectOptions() *expertSelectOptions {
	return &expertSelectOptions{
		status: _defaultExpertStatus,
	}
}

type SelectExpertsOption func(*expertSelectOptions)

func SelectStatus(status entity.Status) SelectExpertsOption {
	return func(eso *expertSelectOptions) {
		eso.status = status
	}
}
