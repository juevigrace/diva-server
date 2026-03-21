package dtos

type SignInDto struct {
	Username    string         `json:"username" validate:"required"`
	Password    string         `json:"password" validate:"required"`
	SessionData SessionDataDto `json:"session_data" validate:"required"`
}

type SignUpDto struct {
	User        CreateUserDto  `json:"user" validate:"required"`
	SessionData SessionDataDto `json:"session_data" validate:"required"`
}

type SessionDataDto struct {
	Device    string `json:"device" validate:"required"`
	IpAddress string `json:"-"`
	UserAgent string `json:"user_agent" validate:"required"`
}

type ForgotPasswordRequestDto struct {
	Email string `json:"email" validate:"required,email"`
}

type ForgotPasswordConfirmDto struct {
	Token       string         `json:"token" validate:"required"`
	SessionData SessionDataDto `json:"session_data" validate:"required"`
}
