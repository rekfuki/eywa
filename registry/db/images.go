package db

import (
	"database/sql"

	"xorm.io/builder"

	"eywa/registry/types"
)

// GetImageWithoutSource returns image from db without source
func (c *Client) GetImageWithoutSource(imageID, userID string) (*types.Image, error) {
	query := c.Builder().
		Select(`i.id, i.user_id, i.registry, i.language,
		 i.name, i.version, i.created_at, i.state, i.size`).
		From("images i").
		Where(builder.Eq{
			"i.id":      imageID,
			"i.user_id": userID,
		})

	var image types.Image
	if err := c.Get(&image, query); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}

	return &image, nil
}

// GetImagesWithoutSource returns images without source
func (c *Client) GetImagesWithoutSource(userID, filter string, pageNumber, perPage int) ([]types.Image, int, error) {
	query := c.Builder().
		Select(`i.id, i.user_id, i.registry, i.language,
		 i.name, i.version, i.created_at, i.state, i.size`).
		From("images i").
		Where(builder.Eq{"i.user_id": userID})

	query = applyImageFilter(query, "i", filter)
	query = query.OrderBy("i.created_at")

	images := []types.Image{}
	total, err := c.SelectWithCount(&images, query, pageNumber, perPage)
	if err != nil {
		return nil, 0, err
	}

	return images, total, nil
}

// CreateImage creates a new image inside the db
func (c *Client) CreateImage(image *types.Image) error {
	query := c.Builder().Insert(builder.Eq{
		"id":         image.ID,
		"user_id":    image.UserID,
		"registry":   image.TaggedRegistry,
		"language":   image.Language,
		"name":       image.Name,
		"version":    image.Version,
		"created_at": image.CreatedAt,
		"state":      image.State,
		"size":       image.Size,
		"source":     image.Source,
	}).Into("images")

	_, err := c.Exec(query)
	return err
}

// UpdateImageState updates the state of an image
func (c *Client) UpdateImageState(imageID, state string) error {
	query := c.Builder().
		Update(builder.Eq{"state": state}).
		From("images").
		Where(builder.Eq{"id": imageID})

	_, err := c.Exec(query)
	return err
}

// DeleteImage deletes image from the db
func (c *Client) DeleteImage(imageID, userID string) error {
	query := c.Builder().
		Delete(builder.Eq{"id": imageID, "user_id": userID}).
		From("images")

	_, err := c.Exec(query)
	return err
}

func applyImageFilter(query *builder.Builder, name string, filter string) *builder.Builder {
	if filter != "" {
		query = query.And(builder.Or(
			ILike{name + ".id::text", filter},
			ILike{name + ".name", filter},
			ILike{name + ".language", filter},
			ILike{name + ".version", filter},
		))
	}

	return query
}
