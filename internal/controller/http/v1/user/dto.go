package userroutes

import "github.com/bogdanshibilov/mindflowbackend/internal/entity"

type userDto struct {
	Id                    string   `json:"id"`
	Email                 string   `json:"email"`
	Name                  string   `json:"name"`
	Roles                 []string `json:"roles"`
	ProfessionalField     string   `json:"professionalField"`
	ExperienceDescription string   `json:"experienceDescription"`
	Phone                 string   `json:"phone"`
}

func userDtoFrom(entity *entity.User) *userDto {
	return &userDto{
		Id:                    entity.Uuid.String(),
		Email:                 entity.Email,
		Name:                  entity.Name,
		Roles:                 entity.Roles,
		ProfessionalField:     entity.ProfessionalField,
		ExperienceDescription: entity.ExperienceDescription,
		Phone:                 entity.Phone,
	}
}

type ForceUpdateUserProfileRequest struct {
	Id                    string `json:"id" binding:"required"`
	Name                  string `json:"name" binding:"required"`
	Email                 string `json:"email" binding:"required"`
	Phone                 string `json:"phone" binding:"required"`
	ProfessionalField     string `json:"professionalField" binding:"required"`
	ExperienceDescription string `json:"experienceDescription" binding:"required"`
}

type DeleteUserByIdRequest struct {
	Id string `json:"id" binding:"required"`
}
