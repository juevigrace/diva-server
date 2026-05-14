package dtos

type VerificationDto struct {
	Token       string          `json:"token" validate:"required,len=6"`
	SessionData *SessionDataDto `json:"session_data" validate:"required"`
}

type RequestVerificationDto struct {
	Email  string `json:"email" validate:"required,email"`
	Action string `json:"action" validate:"required,max=255"`
}
