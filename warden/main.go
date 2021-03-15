package main

import (
	"time"

	"github.com/kelseyhightower/envconfig"
	"github.com/markbates/goth"
	"github.com/markbates/goth/providers/github"
	log "github.com/sirupsen/logrus"

	"eywa/warden/authn"
	"eywa/warden/clients/tugrik"
	"eywa/warden/db"
	"eywa/warden/server"
)

// Config represents wardends startup config
type Config struct {
	DB                     db.Config
	SessionTimeoutDuration time.Duration `envconfig:"warden_session_timeout_duration" default:"15m"`
	SessionSigningKey      string        `envconfig:"warden_session_signing_key" default:"foo-bar"`
	GithubClientID         string        `envconfig:"warden_github_client_id" required:"true"`
	GithubClientSecret     string        `envconfig:"warden_github_client_secret" required:"true"`
	GithubCallbackURL      string        `envconfig:"warden_callback_url" default:"https://tunnel.rekfuki.dev/oauth"`
	TugrikURL              string        `envconfig:"tugrik_url" default:"http://tugrik.faas-system:11080"`
}

func main() {
	var conf Config
	err := envconfig.Process("", &conf)
	if err != nil {
		log.Fatalf("Failed to parse env: %s", err)
	}

	goth.UseProviders(github.New(conf.GithubClientID, conf.GithubClientSecret, conf.GithubCallbackURL, "user:email"))

	migrateDB(conf.DB, 0)

	db, err := db.Connect(conf.DB)
	if err != nil {
		log.Fatalf("Failed to connect to database: %s", err)
	}

	authn.InitTokenCache(db)

	tugrik := tugrik.New(conf.TugrikURL)

	params := &server.ContextParams{
		DB:                     db,
		SessionSigningKey:      conf.SessionSigningKey,
		SessionTimeoutDuration: conf.SessionTimeoutDuration,
		Tugrik:                 tugrik,
	}

	server.Run(params)
}

func migrateDB(dbConf db.Config, target uint) {
	log.Info("Migrating Database Schema")
	db, err := db.Connect(dbConf)
	if err != nil {
		log.Fatalf("failed to connect to trident db: %s", err)
	}

	err = db.Migrate("./migrations", target)
	if err != nil {
		log.Fatalf("failed to create schema: %s", err)
	}

	log.Info("Completed")
}
