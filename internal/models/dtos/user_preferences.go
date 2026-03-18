package dtos

type UserPreferencesDto struct {
	Id                  string `json:"id" validate:"required,uuid"`
	Theme               string `json:"theme" validate:"omitempty,oneof=LIGHT DARK SYSTEM"`
	OnboardingCompleted *bool  `json:"onboarding_completed"`
	Language            string `json:"language" validate:"omitempty,len=2"`
	CreatedAt           int64  `json:"created_at"`
	UpdatedAt           int64  `json:"updated_at"`
}
