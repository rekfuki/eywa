package db

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	log "github.com/sirupsen/logrus"

	// this is needed to enable postgres database support
	_ "github.com/lib/pq"
)

// Client is the client object
type Client struct {
	db *sqlx.DB
	ex sqlx.Ext
}

// Config defines the config information to be passed to Connect method
type Config struct {
	User     string `default:"warden"`
	DBName   string `default:"warden"`
	Password string `envconfig:"warden_db_password" required:"true"`
	Host     string `envconfig:"postgres_host" default:"stolon-proxy.stolon"`
	Port     string `envconfig:"postgres_port" default:"5432"`
}

// NewClient creates a new client object
func NewClient(db *sqlx.DB) *Client {
	return &Client{
		db: db,
		ex: db,
	}
}

// Connect returns db client
func Connect(conf Config) (client *Client, err error) {
	cn := fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s sslmode=disable",
		conf.Host,
		conf.Port,
		conf.DBName,
		conf.User,
		conf.Password)
	rawdb, err := sqlx.Connect("postgres", cn)
	if err != nil {
		log.Errorf("Failed to connect to postgres db: %s", err)
		return
	}
	client = NewClient(rawdb)
	return
}

// DB returns internal sqlx db connection
func (c *Client) DB() *sqlx.DB {
	return c.db
}
