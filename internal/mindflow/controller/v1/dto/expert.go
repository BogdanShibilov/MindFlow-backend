package dto

type BecomeExpertRequest struct {
	ChargePerHour          int    `json:"chargePerHour" binding:"required"`
	ExpertiseAtDescription string `json:"expertiseAtDescription" binding:"required"`
}
