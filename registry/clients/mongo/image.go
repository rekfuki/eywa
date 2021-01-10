package mongo

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"eywa/registry/types"
)

// GetImages returns all the images user has access to
func (c *Client) GetImages(userID string, page, perPage int) (int, []types.Image, error) {
	collection := c.database.Collection(imagesCollection)
	filter := bson.M{"user_id": userID}

	images := []types.Image{}
	total, err := c.FindPaginated(collection, filter, page, perPage, &images)
	if err != nil {
		return 0, nil, err
	}

	return total, images, nil
}

// GetImage returns a specific image without source
func (c *Client) GetImage(id, userID string) (*types.Image, error) {
	return c.getImage(id, userID, false)
}

// GetImageWithSource returns a specific image without source
func (c *Client) GetImageWithSource(id, userID string) (*types.Image, error) {
	return c.getImage(id, userID, true)
}

func (c *Client) getImage(id, userID string, withSource bool) (*types.Image, error) {
	collection := c.database.Collection(imagesCollection)
	filter := bson.M{
		"_id":     id,
		"user_id": userID,
	}

	var image types.Image
	var err error
	if withSource {
		projection := bson.M{"source": 0}
		err = c.FindOneProjection(collection, filter, projection, &image)
	} else {
		err = c.FindOne(collection, filter, &image)
	}
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	return &image, nil
}

// CreateImage creates a new image
func (c *Client) CreateImage(image types.Image) error {
	collection := c.database.Collection(imagesCollection)

	_, err := collection.InsertOne(context.TODO(), image)
	if err != nil {
		return err
	}

	return nil
}

// UpdateImageState updates the state of the image
func (c *Client) UpdateImageState(id, state string) error {
	collection := c.database.Collection(imagesCollection)

	filter := bson.M{"_id": id}
	update := bson.M{
		"$set": bson.M{
			"state": state,
		},
	}
	_, err := collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return err
	}

	return nil
}

// DeleteImage deletes image data
func (c *Client) DeleteImage(id, userID string) error {
	collection := c.database.Collection(imagesCollection)

	filter := bson.M{
		"_id":     id,
		"user_id": userID,
	}
	_, err := collection.DeleteOne(context.TODO(), filter)
	if err != nil {
		return err
	}

	return nil
}
