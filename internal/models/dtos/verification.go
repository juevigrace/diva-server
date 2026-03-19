package dtos

type EmailTokenDto struct {
	Token string `json:"token" validate:"required"`
}
