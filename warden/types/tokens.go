package types

// CreateTokenRequest represents create token request
type CreateTokenRequest struct {
	Name      string `json:"name" binding:"required"`
	ExpiresAt int64  `json:"expires_at"`
}

// AccessTokensResponse represents get all tokens response
type AccessTokensResponse struct {
	Objects []AccessToken `json:"objects"`
	Total   int           `json:"total_count"`
	PerPage int           `json:"per_page"`
	Page    int           `json:"page"`
}

// AccessToken represents access token
type AccessToken struct {
	ID        string `db:"id" json:"id"`
	UserID    string `db:"user_id" json:"-"`
	Name      string `db:"name" json:"name"`
	Token     string `db:"token" json:"token,omitempty"`
	CreatedAt int64  `db:"created_at" json:"created_at"`
	ExpiresAt int64  `db:"expires_at" json:"expires_at,omitempty"`
}
