package dtos

type UserPermissionDto struct {
	UserId       string `json:"user_id" validate:"required,uuid"`
	PermissionId string `json:"permission_id" validate:"required,uuid"`
	Granted      bool   `json:"granted"`
	ExpiresAt    *int64 `json:"expires_at"`
	GrantedBy    string `json:"granted_by" validate:"omitempty,uuid"`
}

type DeleteUserPermissionDto struct {
	UserId       string `json:"user_id" validate:"required,uuid"`
	PermissionId string `json:"permission_id" validate:"required,uuid"`
}
