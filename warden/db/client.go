package db

import (
	"time"

	"github.com/jackc/pgx"
	"github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

// Client is the client object
type Client struct {
	db   *sqlx.DB
	ex   sqlx.Ext
	pool *pgx.ConnPool
}

// Config defines the config information to be passed to Connect method
type Config struct {
	User     string `default:"warden"`
	DBName   string `default:"warden"`
	Password string `envconfig:"warden_db_password" required:"true"`
	Host     string `envconfig:"postgres_host" default:"stolon-proxy.stolon"`
	Port     uint16 `envconfig:"postgres_port" default:"5432"`
}

// Connect returns db client
func Connect(conf Config) (client *Client, err error) {
	connConfig := pgx.ConnConfig{
		Host:     conf.Host,
		Port:     conf.Port,
		Database: conf.DBName,
		User:     conf.User,
		Password: conf.Password,
	}
	connPool, err := pgx.NewConnPool(pgx.ConnPoolConfig{
		ConnConfig:     connConfig,
		AfterConnect:   nil,
		MaxConnections: 20,
		AcquireTimeout: 30 * time.Second,
	})
	if err != nil {
		return nil, errors.Wrap(err, "Call to pgx.NewConnPool failed")
	}

	nativeDB := stdlib.OpenDBFromPool(connPool)
	rawdb, err := sqlx.NewDb(nativeDB, "pgx"), nil
	if err != nil {
		log.Errorf("Failed to connect to postgres db: %s", err)
		return
	}

	client = &Client{
		db:   rawdb,
		ex:   rawdb,
		pool: connPool,
	}

	return
}

// DB returns internal sqlx db connection
func (c *Client) DB() *sqlx.DB {
	return c.db
}
