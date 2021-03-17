package db

import (
	"database/sql"
	"eywa/warden/types"

	"xorm.io/builder"
)

// GetAccessToken returns access token beloging to a user
func (c *Client) GetAccessToken(userID, tokenID string) (*types.AccessToken, error) {
	query := c.Builder().
		Select("*").
		From("access_tokens").
		Where(builder.Eq{"user_id": userID, "id": tokenID})

	var accessToken types.AccessToken
	err := c.Get(&accessToken, query)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &accessToken, nil
}

// GetAccessTokens returns access tokens beloging to a user
func (c *Client) GetAccessTokens(userID, filter string, pageNumber, perPage int) (int, []types.AccessToken, error) {
	query := c.Builder().
		Select("at.id, at.user_id, at.name, at.created_at, at.expires_at").
		From("access_tokens at").
		Where(builder.Eq{"at.user_id": userID})

	query = applyTokenFilter(query, "at", filter)

	accessTokens := []types.AccessToken{}
	total, err := c.SelectWithCount(&accessTokens, query, pageNumber, perPage)
	if err != nil {
		return 0, nil, err
	}

	return total, accessTokens, nil
}

// GetAllAccessTokens returns all access tokens
func (c *Client) GetAllAccessTokens() ([]types.AccessToken, error) {
	query := c.Builder().
		Select("*").
		From("access_tokens")

	accessTokens := []types.AccessToken{}
	if err := c.Select(&accessTokens, query); err != nil {
		return nil, err
	}

	return accessTokens, nil
}

// CreateAccessToken inserts new access token into the database
func (c *Client) CreateAccessToken(accessToken *types.AccessToken) error {
	query := c.Builder().Insert(builder.Eq{
		"id":         accessToken.ID,
		"user_id":    accessToken.UserID,
		"name":       accessToken.Name,
		"token":      accessToken.Token,
		"created_at": accessToken.CreatedAt,
		"expires_at": accessToken.ExpiresAt,
	}).Into("access_tokens")

	_, err := c.Exec(query)
	if err != nil {
		return err
	}

	return nil
}

// DeleteAccessToken deletes access token from the database
func (c *Client) DeleteAccessToken(userID, tokenID string) error {
	query := c.Builder().
		Delete(builder.Eq{"user_id": userID, "id": tokenID}).
		From("access_tokens")

	_, err := c.Exec(query)
	if err != nil {
		return err
	}

	return nil
}

func applyTokenFilter(query *builder.Builder, name string, filter string) *builder.Builder {
	if filter != "" {
		query = query.And(builder.Or(
			ILike{name + ".id::text", filter},
			ILike{name + ".name::text", filter},
		))
	}

	return query
}
