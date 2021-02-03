package types

import (
	"time"
)

// User is a structure representing user data
type User struct {
	ID                 string    `json:"id" db:"id"`
	Name               string    `json:"name" db:"name"`
	AvatarURL          string    `json:"avatar_url" db:"avatar_url"`
	OauthProvider      string    `json:"oauth_provider" db:"oauth_provider"`
	OauthProviderID    string    `json:"provider_id" db:"oauth_provider_id"`
	OauthProviderEmail string    `json:"oauth_provider_email" db:"oauth_provider_email"`
	OauthProviderLogin string    `json:"oauth_provider_login" db:"oauth_provider_login"`
	CreatedAt          time.Time `json:"created_at" db:"created_at" diff:"-"`
	LastSeenAt         time.Time `json:"last_seen_at" db:"last_seen_at" diff:"-"`
}
