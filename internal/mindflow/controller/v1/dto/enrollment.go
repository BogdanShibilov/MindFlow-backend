package dto

type CreateEnrollmentRequest struct {
	MentorId        string `json:"mentorId" binding:"required"`
	MenteeQuestions string `json:"menteeQuestions" binding:"required"`
}

type EnrollmentsByMemberIdRequest struct {
	ByWhoseId string `json:"byWhoseId" binding:"required"`
	Id        string `json:"id" binding:"required"`
}
