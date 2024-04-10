package dto

type CreateExpertRequest struct {
	Uuid                   string `json:"uuid" binding:"required"`
	Position               string `json:"position" binding:"required"`
	ChargePerHour          int    `json:"chargePerHour" binding:"required"`
	ExperienceDescription  string `json:"experienceDescription" binding:"required"`
	ExpertiseAtDescription string `json:"expertiseAtDescription" binding:"required"`
}
