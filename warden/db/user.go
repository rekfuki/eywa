package db

import (
	"database/sql"

	"xorm.io/builder"

	"eywa/warden/types"
)

// GetUserByUserID returns user by user id
func (c *Client) GetUserByUserID(userID string) (*types.User, error) {
	query := c.Builder().
		Select("u.*").
		From("users u").
		Where(builder.Eq{"u.oauth_provider_id": userID})

	var user types.User
	err := c.Get(&user, query)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}

// CreateUser creates a new user in the database
func (c *Client) CreateUser(user *types.User) (*types.User, error) {
	query := c.Builder().
		Insert(builder.Eq{
			"id":                   user.ID,
			"name":                 user.Name,
			"avatar_url":           user.AvatarURL,
			"oauth_provider":       user.OauthProvider,
			"oauth_provider_id":    user.OauthProviderID,
			"oauth_provider_email": user.OauthProviderEmail,
			"oauth_provider_login": user.OauthProviderLogin,
			"created_at":           user.CreatedAt,
			"last_seen_at":         user.LastSeenAt,
		}).
		Into("users")

	_, err := c.Exec(query)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// UpdateUser updates a user in the database
func (c *Client) UpdateUser(user *types.User) error {
	query := c.Builder().
		Update(builder.Eq{
			"name":                 user.Name,
			"avatar_url":           user.AvatarURL,
			"oauth_provider":       user.OauthProvider,
			"oauth_provider_email": user.OauthProviderEmail,
			"oauth_provider_login": user.OauthProviderLogin,
			"last_seen_at":         user.LastSeenAt,
		}).
		From("users").
		Where(builder.Eq{"id": user.ID})

	_, err := c.Exec(query)
	if err != nil {
		return err
	}

	return nil
}
