package responses

type UserPermissionResponse struct {
	PermissionId string  `json:"permission_id"`
	GrantedBy    *string `json:"granted_by"`
	Granted      bool    `json:"granted"`
	GrantedAt    *int64  `json:"granted_at"`
	ExpiresAt    *int64  `json:"expires_at"`
	UpdatedAt    int64   `json:"updated_at"`
}
