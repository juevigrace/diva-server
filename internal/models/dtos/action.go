package dtos

type UserActionDto struct {
	Action string `json:"action" validate:"required"`
}
