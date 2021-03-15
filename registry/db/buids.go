package db

import (
	"database/sql"

	"github.com/lib/pq"
	"xorm.io/builder"

	"eywa/registry/types"
)

// GetBuild returns a build form the db
func (c *Client) GetBuild(imageID, userID string) (*types.Build, error) {
	query := c.Builder().Select("*").
		From("builds").
		Where(builder.Eq{"image_id": imageID, "user_id": userID})

	var build types.Build
	err := c.Get(&build, query)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}

	return &build, err
}

// CreateBuild inserts a new build into db
func (c *Client) CreateBuild(build *types.Build) error {
	query := c.Builder().Insert(builder.Eq{
		"image_id":   build.ImageID,
		"user_id":    build.UserID,
		"logs":       build.Logs,
		"state":      build.State,
		"created_at": build.CreatedAt,
	}).Into("builds")

	_, err := c.Exec(query)
	return err
}

// UpdateBuild updates a build inside db
func (c *Client) UpdateBuild(imageID, state string, logs pq.StringArray) error {
	query := c.Builder().
		Update(builder.Eq{
			"logs":  logs,
			"state": state,
		}).
		From("builds").
		Where(builder.Eq{"image_id": imageID})

	_, err := c.Exec(query)
	return err
}

// DeleteBuild deletes build info form the db
func (c *Client) DeleteBuild(imageID, userID string) error {
	query := c.Builder().
		Delete(builder.Eq{"image_id": imageID, "user_id": userID}).
		From("builds")

	_, err := c.Exec(query)
	return err
}
