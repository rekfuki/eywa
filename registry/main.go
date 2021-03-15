package main

import (
	"flag"

	"github.com/kelseyhightower/envconfig"
	log "github.com/sirupsen/logrus"

	"eywa/registry/builder"
	"eywa/registry/clients/docker"
	"eywa/registry/db"
	"eywa/registry/server"
)

// Config represents registry startup configuration
type Config struct {
	Postgres         db.Config
	Registry         string `envconfig:"registry_url" default:"registry.eywa.rekfuki.dev"`
	RegistryUser     string `envconfig:"registry_user" required:"true"`
	RegistryPassword string `envconfig:"registry_password" required:"true"`
	NumWorkers       int    `envconfig:"builder_worker_count" default:"3"`
}

func main() {
	var conf Config
	err := envconfig.Process("", &conf)
	if err != nil {
		log.Fatalf("Failed to parse env: %s", err)
	}

	debug := flag.Bool("debug", false, "(optional) set log level to debug")
	flag.Parse()

	if *debug {
		log.SetLevel(log.DebugLevel)
	}

	migrateDB(conf.Postgres, 0)

	db, err := db.Connect(conf.Postgres)
	if err != nil {
		log.Fatalf("failed to connect to registry db: %s", err)
	}

	docker := docker.New(conf.Registry, conf.RegistryUser, conf.RegistryPassword)

	builder, err := builder.New(&builder.Config{
		Registry:         conf.Registry,
		RegistryUser:     conf.RegistryUser,
		RegistryPassword: conf.RegistryPassword,
		DB:               db,
		NumWorkers:       3,
	})
	if err != nil {
		log.Fatalf("Failed to setup builder: %s", err)
	}
	go builder.Start()

	params := &server.ContextParams{
		DB:      db,
		Builder: builder,
		Docker:  docker,
	}

	server.Run(params)

	log.Exit(0)
}

func migrateDB(dbConf db.Config, target uint) {
	log.Info("Migrating Database Schema")
	db, err := db.Connect(dbConf)
	if err != nil {
		log.Fatalf("failed to connect to registry db: %s", err)
	}

	err = db.Migrate("./migrations", target)
	if err != nil {
		log.Fatalf("failed to create schema: %s", err)
	}

	log.Info("Completed")
}
