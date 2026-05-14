package dtos

type CreateUserDto struct {
	Email    string `json:"email" validate:"required,email,max=100"`
	Username string `json:"username" validate:"required,min=3,max=50"`
	Password string `json:"password" validate:"required,min=4"`
}

type CreateProfileDto struct {
	FirstName string `json:"first_name" validate:"required,max=255"`
	LastName  string `json:"last_name" validate:"required,max=255"`
	Alias     string `json:"alias" validate:"required,max=255"`
	Bio       string `json:"bio" validate:"omitempty,max=255"`
	BirthDate int64  `json:"birth_date" validate:"required,gt=0"`
}

type UpdateProfileDto struct {
	ID        string `json:"user_id" validate:"required,uuid"`
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
	NewPassword string `json:"new_password" validate:"required,min=4"`
}

type UpdatePhoneNumberDto struct {
	PhoneNumber string `json:"phone_number" validate:"required"`
}

type UpdateEmailDto struct {
	Email string `json:"email" validate:"required,email"`
}

type UserPermissionDto struct {
	UserId       string `json:"user_id" validate:"required,uuid"`
	PermissionId string `json:"permission_id" validate:"required,uuid"`
	Granted      bool   `json:"granted" validate:"required"`
	ExpiresAt    *int64 `json:"expires_at" validate:"required"`
}

type DeleteUserPermissionDto struct {
	UserId       string `json:"user_id" validate:"required,uuid"`
	PermissionId string `json:"permission_id" validate:"required,uuid"`
}

type UserPreferencesDto struct {
	Id                  string `json:"id" validate:"required,uuid"`
	Device              string `json:"device" validate:"required,max=100"`
	Theme               string `json:"theme" validate:"required,oneof=LIGHT DARK SYSTEM"`
	OnboardingCompleted bool   `json:"onboarding_completed" validate:"required"`
	Language            string `json:"language" validate:"required,max=10"`
	CreatedAt           int64  `json:"created_at" validate:"required,gt=0"`
	UpdatedAt           int64  `json:"updated_at"`
}
