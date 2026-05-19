package responses

type UserResponse struct {
	ID          string `json:"id"`
	Username    string `json:"username"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phone_number"`
	Verified    bool   `json:"verified"`
	Role        string `json:"role"`
	CreatedAt   int64  `json:"created_at"`
	UpdatedAt   int64  `json:"updated_at"`
	DeletedAt   *int64 `json:"deleted_at"`
}

type UserProfileResponse struct {
	UserID      *string `json:"user_id,omitempty"`
	FirstName   string  `json:"first_name"`
	LastName    string  `json:"last_name"`
	BirthDate   int64   `json:"birth_date"`
	PhoneNumber string  `json:"phone_number"`
	Alias       string  `json:"alias"`
	Avatar      string  `json:"avatar"`
	Bio         string  `json:"bio"`
}

type UserPermissionResponse struct {
	UserID       *string `json:"user_id,omitempty"`
	PermissionID string  `json:"permission_id"`
	GrantedBy    *string `json:"granted_by"`
	Granted      bool    `json:"granted"`
	GrantedAt    *int64  `json:"granted_at"`
	ExpiresAt    *int64  `json:"expires_at"`
	UpdatedAt    int64   `json:"updated_at"`
}

type UserPreferencesResponse struct {
	UserID              *string `json:"user_id,omitempty"`
	Id                  string  `json:"id"`
	Theme               string  `json:"theme"`
	OnboardingCompleted bool    `json:"onboarding_completed"`
	Language            string  `json:"language"`
	LastSyncAt          int64   `json:"last_sync_at"`
	CreatedAt           int64   `json:"created_at"`
	UpdatedAt           int64   `json:"updated_at"`
}

type UserActionResponse struct {
	UserID     *string `json:"user_id,omitempty"`
	ID         string  `json:"id"`
	ActionName string  `json:"action_name"`
}
