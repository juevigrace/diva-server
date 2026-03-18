package responses

type SessionResponse struct {
	SessionId    string `json:"session_id"`
	UserId       string `json:"user_id"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	Status       string `json:"status"`
	Device       string `json:"device"`
	Ip           string `json:"ip"`
	Agent        string `json:"agent"`
	ExpiresAt    int64  `json:"expires_at"`
	CreatedAt    int64  `json:"created_at"`
	UpdatedAt    int64  `json:"updated_at"`
}
