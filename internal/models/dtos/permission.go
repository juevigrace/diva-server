package dtos

type CreatePermissionDto struct {
	Name        string `json:"name" validate:"required,max=255"`
	Description string `json:"description" validate:"required,max=255"`
	Action      string `json:"action" validate:"required,max=255"`
	RoleLevel   string `json:"level" validate:"required,oneof=USER MODERATOR ADMIN,max=20"`
}

type UpdatePermissionDto struct {
	ID          string `json:"id" validate:"required,uuid"`
	Name        string `json:"name" validate:"required,max=255"`
	Description string `json:"description" validate:"required,max=255"`
}

type UpdatePermissionRoleLevelDto struct {
	ID        string `json:"id" validate:"required,uuid"`
	RoleLevel string `json:"level" validate:"required,oneof=USER MODERATOR ADMIN,max=20"`
}
