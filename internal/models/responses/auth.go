package responses

type AuthResponse struct {
	Session *SessionResponse  `json:"session"`
	Actions []*ActionResponse `json:"actions"`
}
