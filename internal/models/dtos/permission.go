package dtos

type UpdatePermissionDto struct {
	Name        string `json:"name" validate:"required,max=255"`
	Description string `json:"description" validate:"required,max=255"`
}

type UpdatePermissionRoleLevelDto struct {
	Level string `json:"level" validate:"required,oneof=USER MODERATOR ADMIN,max=20"`
}
