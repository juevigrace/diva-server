package responses

type UserPreferencesResponse struct {
	Id                  string `json:"id"`
	Theme               string `json:"theme"`
	OnboardingCompleted bool   `json:"onboarding_completed"`
	Language            string `json:"language"`
	LastSyncAt          int64  `json:"last_sync_at"`
	CreatedAt           int64  `json:"created_at"`
	UpdatedAt           int64  `json:"updated_at"`
}
