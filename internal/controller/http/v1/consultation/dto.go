package consultationroute

import "time"

type applyForConsultationRequest struct {
	ExpertId        string `json:"expertId" binding:"required"`
	MenteeQuestions string `json:"menteeQuestions" binding:"required"`
}

type createMeetingRequest struct {
	ConsultationId string    `json:"consultationId"`
	StartTime      time.Time `json:"startTime"`
	Link           string    `json:"link"`
}
