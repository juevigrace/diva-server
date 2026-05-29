package dtos

type CreateUserDto struct {
	Email    string `json:"email" validate:"required,email,max=100"`
	Username string `json:"username" validate:"required,min=3,max=50"`
	Password string `json:"password" validate:"required,min=4,max=255"`
}

type CreateProfileDto struct {
	FirstName string `json:"first_name" validate:"required,max=255"`
	LastName  string `json:"last_name" validate:"required,max=255"`
	Alias     string `json:"alias" validate:"required,max=255"`
	Bio       string `json:"bio" validate:"omitempty,max=255"`
	BirthDate int64  `json:"birth_date" validate:"required,gt=0"`
}

type UpdateProfileDto struct {
	FirstName string `json:"first_name" validate:"required,max=255"`
	LastName  string `json:"last_name" validate:"required,max=255"`
	Alias     string `json:"alias" validate:"required,max=255"`
	Bio       string `json:"bio" validate:"omitempty,max=255"`
	BirthDate int64  `json:"birth_date" validate:"required,gt=0"`
}

type UpdateUsernameDto struct {
	Username string `json:"username" validate:"required,min=3,max=50"`
}

type UpdatePasswordDto struct {
	NewPassword string `json:"new_password" validate:"required,min=4,max=255"`
}

type UpdatePhoneNumberDto struct {
	PhoneNumber string `json:"phone_number" validate:"required,max=30"`
}

type UpdateEmailDto struct {
	Email string `json:"email" validate:"required,email,max=100"`
}

type UpdateRole struct {
	Role string `json:"role" validate:"required,oneof=USER MODERATOR ADMIN"`
}

type UpdateVerified struct {
	Verified bool `json:"verified" validate:"required"`
}

type UserPermissionDto struct {
	UserId       string `json:"user_id" validate:"required,uuid"`
	PermissionId string `json:"permission_id" validate:"required,uuid"`
	Granted      bool   `json:"granted" validate:"required"`
	ExpiresAt    *int64 `json:"expires_at" validate:"required,omitempty,gt=0"`
}

type CreateUserPreferencesDto struct {
	Device              string `json:"-"`
	Theme               string `json:"theme" validate:"required,oneof=LIGHT DARK SYSTEM"`
	OnboardingCompleted bool   `json:"onboarding_completed" validate:"required"`
	Language            string `json:"language" validate:"required,max=10"`
}

type UpdateUserPreferencesDto struct {
	Theme    string `json:"theme" validate:"required,oneof=LIGHT DARK SYSTEM"`
	Language string `json:"language" validate:"required,max=10"`
}
