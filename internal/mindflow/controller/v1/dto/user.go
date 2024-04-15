package dto

type UpdateUserDetailsRequest struct {
	UserId                string `json:"userId" binding:"required"`
	PhoneNumber           string `json:"phoneNumber" binding:"required,e164"`
	ProfessionalField     string `json:"professionalField" binding:"required"`
	ExperienceDescription string `json:"experienceDescription" binding:"required"`
}
