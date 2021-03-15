package db

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"eywa/tugrik/types"
)

// GetUserDatabase returns user database information
func (c *Client) GetUserDatabase(userID string) (*types.UserDatabaseInfo, error) {
	ctx := context.Background()

	userDB := c.client.Database(userID)
	dbStats := userDB.RunCommand(ctx, bson.D{primitive.E{Key: "dbStats", Value: 1}})
	if dbStats.Err() != nil {
		return nil, dbStats.Err()
	}

	var userDBStats types.UserDatabaseInfo
	if err := dbStats.Decode(&userDBStats); err != nil {
		return nil, err
	}

	collectionNames, err := userDB.ListCollectionNames(ctx, bson.M{})
	if err != nil {
		return nil, err
	}

	collections := []types.CollectionInfo{}
	for _, col := range collectionNames {
		result := userDB.RunCommand(ctx, bson.D{primitive.E{Key: "collStats", Value: col}})
		if result.Err() != nil {
			return nil, result.Err()
		}

		var out types.CollectionInfo
		if err := result.Decode(&out); err != nil {
			return nil, err
		}

		collections = append(collections, out)
	}

	userDBStats.CollectionsInfo = collections
	return &userDBStats, nil
}

// CheckUserExists checks if the user already exists in the database1
func (c *Client) CheckUserExists(userID string) (bool, error) {
	result := c.database.Collection("system.users").
		FindOne(context.Background(), bson.M{"user": userID})

	err := result.Err()
	if err != nil {
		if c.IsNotFound(err) {
			return false, nil
		}

		return false, err
	}

	return true, nil
}

// CreateUserDatabase cretes a new user adn gives them access to their db in mongo
func (c *Client) CreateUserDatabase(userID, password string) error {
	result := c.client.Database(userID).RunCommand(
		context.Background(),
		bson.D{
			primitive.E{Key: "createUser", Value: userID},
			primitive.E{Key: "pwd", Value: password},
			primitive.E{Key: "roles", Value: []string{"readWrite"}},
		},
	)

	return result.Err()
}

// DropCollection drops a collection from the user database
func (c *Client) DropCollection(userID, collectionName string) error {
	return c.client.Database(userID).
		Collection(collectionName).
		Drop(context.Background())
}
