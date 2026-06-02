package dtos

type VerifyActionDto struct {
	ActionID string `json:"action_id" validate:"required,uuid"`
	Token    string `json:"token" validate:"required,len=6"`
}

type RequestActionVerificationDto struct {
	Email  string `json:"email" validate:"required,email,max=100"`
	Action string `json:"action" validate:"required,max=255"`
}
