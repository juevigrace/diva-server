package dtos

type SignInDto struct {
	Username    string         `json:"username" validate:"required,max=100"`
	Password    string         `json:"password" validate:"required,max=1000"`
	SessionData SessionDataDto `json:"session_data" validate:"required"`
}

type SignUpDto struct {
	User        CreateUserDto  `json:"user" validate:"required"`
	SessionData SessionDataDto `json:"session_data" validate:"required"`
}

type SessionDataDto struct {
	Device    string `json:"device" validate:"required,max=100"`
	IpAddress string `json:"-"`
	UserAgent string `json:"user_agent" validate:"required,max=255"`
}

type ForgotPasswordConfirmDto struct {
	ActionID    string         `json:"id" validate:"required,uuid"`
	SessionData SessionDataDto `json:"session_data" validate:"required"`
}
