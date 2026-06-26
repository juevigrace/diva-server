package responses

// TODO: i might need to change this because of the expiration times
type SessionResponse struct {
	SessionId        string `json:"session_id"`
	UserId           string `json:"user_id"`
	AccessToken      string `json:"access_token"`
	RefreshToken     string `json:"refresh_token"`
	Status           string `json:"status"`
	Type             string `json:"type"`
	Device           string `json:"device"`
	Ip               string `json:"ip"`
	Agent            string `json:"agent"`
	AccessExpiresAt  int64  `json:"access_expires_at"`
	RefreshExpiresAt int64  `json:"refresh_expires_at"`
	CreatedAt        int64  `json:"created_at"`
	UpdatedAt        int64  `json:"updated_at"`
}
