package responses

type UserResponse struct {
	ID           string `json:"id"`
	Email        string `json:"email"`
	Username     string `json:"username"`
	BirthDate    int64  `json:"birth_date"`
	PhoneNumber  string `json:"phone_number"`
	Alias        string `json:"alias"`
	Avatar       string `json:"avatar"`
	Bio          string `json:"bio"`
	UserVerified bool   `json:"user_verified"`
	CreatedAt    int64  `json:"created_at"`
	UpdatedAt    int64  `json:"updated_at"`
	DeletedAt    *int64 `json:"deleted_at"`
}
