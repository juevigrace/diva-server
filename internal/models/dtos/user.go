package dtos

type CreateUserDto struct {
	Email    string `json:"email" validate:"required,email"`
	Username string `json:"username" validate:"required,min=3,max=50"`
	Password string `json:"password" validate:"required,min=4"`
	Alias    string `json:"alias" validate:"omitempty,max=50"`
}

type UpdateProfileDto struct {
	Alias     string `json:"alias" validate:"required"`
	Bio       string `json:"bio" validate:"omitempty,max=500"`
	Avatar    string `json:"avatar" validate:"omitempty,url"`
	BirthDate int64  `json:"birth_date" validate:"omitempty"`
}

type UpdateUsernameDto struct {
	Username string `json:"username" validate:"required,min=3,max=50"`
}

type UpdatePhoneNumberDto struct {
	PhoneNumber string `json:"phone_number" validate:"required"`
}

type UpdateEmailDto struct {
	Email string `json:"email" validate:"required,email"`
}
