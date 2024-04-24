package authroutes

type signUpRequest struct {
	Username              string `json:"username" binding:"required"`
	Password              string `json:"password" binding:"required,min=5"`
	Email                 string `json:"email" binding:"required,email"`
	Phone                 string `json:"phone" binding:"required,e164"`
	ProfessionalField     string `json:"professionalField" binding:"required"`
	ExperienceDescription string `json:"experienceDescription" binding:"required"`
}

type signInWithEmailRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=5"`
}

type signInWithEmailResponse struct {
	AccessToken string `json:"accessToken"`
}
