package expertroutes

import (
	"github.com/bogdanshibilov/mindflowbackend/internal/entity"
)

type applyForExpertRequest struct {
	HelpDescription string `json:"helpDescription" binding:"required"`
	Price           int    `json:"price" binding:"required"`
}

type changeStatusExpertRequest struct {
	ExpertId string `json:"expertId" binding:"required"`
	Status   int    `json:"status" binding:"required"`
}

type expertDTO struct {
	UserId                string `json:"userId"`
	Email                 string `json:"email"`
	Name                  string `json:"name"`
	Phone                 string `json:"phone"`
	ProfessionalField     string `json:"professionalField"`
	ExperienceDescription string `json:"experienceDescription"`
	HelpDescription       string `json:"helpDescription"`
	Price                 int    `json:"price"`
}

func expertDtoFrom(entity *entity.Expert) *expertDTO {
	return &expertDTO{
		UserId:                entity.UserUuid.String(),
		Email:                 entity.Email,
		Name:                  entity.Name,
		Phone:                 entity.Phone,
		ProfessionalField:     entity.ProfessionalField,
		ExperienceDescription: entity.ExperienceDescription,
		HelpDescription:       entity.HelpDescription,
		Price:                 entity.Price,
	}
}
