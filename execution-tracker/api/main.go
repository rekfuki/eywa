package main

import (
	"github.com/kelseyhightower/envconfig"
	log "github.com/sirupsen/logrus"

	"eywa/execution-tracker/api/server"
	"eywa/execution-tracker/db"
)

func main() {
	var conf server.Config
	err := envconfig.Process("", &conf)
	if err != nil {
		log.Fatalf("Failed to parse env: %s", err)
	}

	migrateDB(conf.Postgres, 0)

	db, err := db.Connect(conf.Postgres)
	if err != nil {
		log.Fatalf("failed to connect to sawmill db: %s", err)
	}

	server.Run(conf, db)
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
