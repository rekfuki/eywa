package db

import (
	"context"
	"time"

	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Client is the client object
type Client struct {
	client   *mongo.Client
	database *mongo.Database
}

// Config defines the config information to be passed to Connect method
type Config struct {
	Host     string `envconfig:"mongodb_host" default:"mongodb.mongodb:27017"`
	Database string `envconfig:"-" default:"admin"`
	User     string `envconfig:"-" default:"root"`
	Password string `envconfig:"tugrik_db_password"`
}

const (
	defaultMaxQueryTime = 2 * time.Second
)

// Connect connects to a MongoDB cluster and returns a client
func Connect(conf Config) (*Client, error) {
	log.Infof("Connecting to MongoDB @ %s", conf.Host)

	opts := options.Client().
		SetHosts([]string{conf.Host}).
		SetAuth(options.Credential{
			Username:      conf.User,
			Password:      conf.Password,
			AuthSource:    conf.Database,
			AuthMechanism: "SCRAM-SHA-1",
		})
	opts = opts.SetConnectTimeout(time.Second * 10)
	client, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		return nil, err
	}
	database := client.Database(conf.Database)
	return &Client{
		client:   client,
		database: database,
	}, nil
}

// IsNotFound returns whether err informs of documents not matching the search filter
func (c Client) IsNotFound(err error) bool {
	return err == mongo.ErrNoDocuments
}

// IsConflict returns whether err indicates a conflict
func (c Client) IsConflict(err error) bool {
	return mongo.IsDuplicateKeyError(err)
}

// Find is a wrapper around the regular mongo driver Find function that adds a query execution timeout
func (c *Client) Find(collection *mongo.Collection, filter interface{}, result interface{}) error {
	ctx := context.TODO()
	opts := options.Find().SetMaxTime(defaultMaxQueryTime)
	cursor, err := collection.Find(ctx, filter, opts)
	return c.readAll(ctx, cursor, err, result)
}

func (c *Client) readAll(ctx context.Context, cursor *mongo.Cursor, err error, result interface{}) error {
	if err != nil {
		if c.IsNotFound(err) {
			return nil
		}
		log.Errorf("DatabaseError: %s", err)
		return err
	}
	defer cursor.Close(ctx)
	err = cursor.All(ctx, result)
	if err != nil {
		if c.IsNotFound(err) {
			return nil
		}
		log.Errorf("DatabaseError: %s", err)
	}
	return err
}

// FindOne is a wrapper around the regular mongo driver FindOne function that adds a query execution timeout
func (c *Client) FindOne(collection *mongo.Collection, filter interface{}, result interface{}) error {
	opts := options.FindOne().SetMaxTime(defaultMaxQueryTime)
	return collection.FindOne(context.TODO(), filter, opts).Decode(result)
}

// Aggregate is a wrapper around the regular mongo driver Aggregate function that adds a query execution timeout
func (c *Client) Aggregate(collection *mongo.Collection, pipeline mongo.Pipeline, result interface{}) error {
	opts := options.Aggregate()
	ctx := context.TODO()
	cursor, err := collection.Aggregate(ctx, pipeline, opts)
	return c.readAll(ctx, cursor, err, result)
}

// UpdateOne is a wrapper around the regular mongo driver UpdateOne function that would add a query execution timeout if it were possible
func (c *Client) UpdateOne(collection *mongo.Collection, filter interface{}, update interface{}) error {
	opts := options.Update()
	_, err := collection.UpdateOne(context.TODO(), filter, update, opts)
	return err
}

// Database gives sneaky access to the database field
// TODO This is just to support temporary tdb - remove this when that's replaced with trident's
func (c *Client) Database() *mongo.Database {
	return c.database
}
