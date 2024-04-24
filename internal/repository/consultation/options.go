package consultationrepo

type WhoseUuid string

const (
	ByExpertUuid WhoseUuid = "expert_uuid"
	ByMenteeUuid WhoseUuid = "mentee_uuid"
)

const (
	_defaultWhoseUuid WhoseUuid = ByExpertUuid
)

type selectByPersonUuidOptions struct {
	whoseUuid WhoseUuid
}

func newDefaultSelectByPersonUuidOptions() *selectByPersonUuidOptions {
	return &selectByPersonUuidOptions{
		whoseUuid: _defaultWhoseUuid,
	}
}

type ByPersonUuidOption func(*selectByPersonUuidOptions)

func SelectByWhoseUuid(whoseUuid WhoseUuid) ByPersonUuidOption {
	return func(sbpuo *selectByPersonUuidOptions) {
		sbpuo.whoseUuid = whoseUuid
	}
}
