package consultationroute

type applyForConsultationRequest struct {
	ExpertId        string `json:"expertId" binding:"required"`
	MenteeQuestions string `json:"menteeQuestions" binding:"required"`
}
