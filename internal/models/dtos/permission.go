package dtos

type CreatePermissionDto struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description" validate:"required"`
	Action      string `json:"action" validate:"required"`
	RoleLevel   string `json:"level" validate:"required,oneof=USER MODERATOR ADMIN"`
}

type UpdatePermissionDto struct {
	ID          string `json:"id" validate:"required"`
	Name        string `json:"name" validate:"required"`
	Description string `json:"description" validate:"required"`
}

type UpdatePermissionRoleLevelDto struct {
	ID        string `json:"id" validate:"required"`
	RoleLevel string `json:"level" validate:"required,oneof=USER MODERATOR ADMIN"`
}
