package dtos

type CreatePermissionDto struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description" validate:"required"`
	RoleLevel   string `json:"level" validate:"required,oneof=USER MODERATOR ADMIN"`
}

type UpdatePermissionDto struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description" validate:"required"`
}
