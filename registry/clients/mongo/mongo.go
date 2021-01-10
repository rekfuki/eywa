package mongo

import (
	"context"
	"time"

	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Config defines the config information to be passed to Connect method
type Config struct {
	Hosts    []string `envconfig:"mongo_host" default:"0.0.0.0:27017"`
	Database string   `envconfig:"-" default:"registry"`
	User     string   `envconfig:"-" default:"registry"`
	Password string   `envconfig:"registry_mongo_password" default:"foobar"`
}

const (
	imagesCollection   = "images"
	servicesCollection = "services"

	defaultMaxQueryTime = 2 * time.Second
)

// Client is the client object
type Client struct {
	client   *mongo.Client
	database *mongo.Database
	db       string
}

// Connect connects to a MongoDB cluster and returns a client
func Connect(conf Config) (*Client, error) {
	log.Infof("Connecting to MongoDB @ %s", conf.Hosts)
	opts := options.Client().
		SetHosts(conf.Hosts).
		SetAuth(options.Credential{
			Username:   conf.User,
			Password:   conf.Password,
			AuthSource: conf.Database,
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
		db:       conf.Database,
	}, nil
}

// Aggregate is a wrapper around the regular mongo driver Aggregate function that adds a query execution timeout
func (c *Client) Aggregate(collection *mongo.Collection, pipeline mongo.Pipeline, result interface{}) error {
	opts := options.Aggregate()
	ctx := context.TODO()
	cursor, err := collection.Aggregate(ctx, pipeline, opts)
	return c.readAll(ctx, cursor, err, result)
}

// FindOpts is a wrapper around the regular mongo driver find function that adds a query execution timeout
func FindOpts() *options.FindOptions {
	return options.Find().
		SetMaxTime(defaultMaxQueryTime)
}

// FindOne is a wrapper around the regular mongo driver FindOne function that adds a query execution timeout
func (c *Client) FindOne(collection *mongo.Collection, criteria interface{}, result interface{}) error {
	return collection.FindOne(context.TODO(), criteria, options.FindOne().
		SetMaxTime(defaultMaxQueryTime)).Decode(result)
}

// FindOneProjection is a wrapper around the regular mongo driver FindOne function that adds a query execution timeout
func (c *Client) FindOneProjection(collection *mongo.Collection, criteria interface{}, projection interface{}, result interface{}) error {
	return collection.FindOne(context.TODO(), criteria, options.FindOne().SetProjection(projection).
		SetMaxTime(defaultMaxQueryTime)).Decode(result)
}

// Find is a wrapper around the regular mongo driver Find function that adds a query execution timeout
func (c *Client) Find(collection *mongo.Collection, filter interface{}, result interface{}) error {
	ctx := context.TODO()
	opts := options.Find().SetMaxTime(defaultMaxQueryTime)
	cursor, err := collection.Find(ctx, filter, opts)
	return c.readAll(ctx, cursor, err, result)
}

// FindPaginated is a warapper around combined Find and Count
func (c *Client) FindPaginated(collection *mongo.Collection, filter interface{}, page, perPage int, result interface{}) (int, error) {
	ctx := context.TODO()
	total, err := collection.CountDocuments(ctx, filter, options.Count())
	if err != nil {
		return 0, err
	}

	opts := options.Find().SetMaxTime(defaultMaxQueryTime)
	if perPage > 0 && page > 0 {
		opts = opts.
			SetSkip(int64(perPage * (page - 1))).
			SetLimit(int64(perPage))
	}
	cursor, err := collection.Find(ctx, filter, opts)
	return int(total), c.readAll(ctx, cursor, err, result)
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

// IsNotFound returns whether err informs of documents not matching the search filter
func (c *Client) IsNotFound(err error) bool {
	return err == mongo.ErrNoDocuments
}
