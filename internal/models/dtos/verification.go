package dtos

type VerificationDto struct {
	Token       string          `json:"token" validate:"required"`
	SessionData *SessionDataDto `json:"session_data"`
}

type RequestVerificationDto struct {
	Email  string `json:"email" validate:"required"`
	Action string `json:"action" validate:"required"`
}
