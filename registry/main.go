package main

import (
	"flag"

	"github.com/kelseyhightower/envconfig"
	log "github.com/sirupsen/logrus"

	"eywa/registry/builder"
	"eywa/registry/clients/docker"
	"eywa/registry/clients/mongo"
	"eywa/registry/server"
)

// Config represents registry startup configuration
type Config struct {
	Mongo            mongo.Config
	Registry         string `envconfig:"registry_url" default:"registry.eywa.rekfuki.dev"`
	RegistryUser     string `envconfig:"registry_user" required:"true"`
	RegistryPassword string `envconfig:"registry_password" required:"true"`
	NumWorkers       int    `envconfig:"builder_worker_count" default:"3"`
	// GatewayURL       string `envconfig:"gateway_url" default:"gateway.faas-system:8080"`
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

	mc, err := mongo.Connect(conf.Mongo)
	if err != nil {
		log.Fatalf("Failed to connecto to mongo: %s", err)
	}

	docker := docker.New(conf.Registry, conf.RegistryUser, conf.RegistryPassword)
	// gateway := gateway.New(conf.GatewayURL)

	builder, err := builder.New(&builder.Config{
		Registry:         conf.Registry,
		RegistryUser:     conf.RegistryUser,
		RegistryPassword: conf.RegistryPassword,
		Mongo:            mc,
		NumWorkers:       3,
	})
	if err != nil {
		log.Fatalf("Failed to setup builder: %s", err)
	}
	go builder.Start()

	params := &server.ContextParams{
		Mongo:   mc,
		Builder: builder,
		Docker:  docker,
		// Gateway: gateway,
	}

	server.Run(params)

	log.Exit(0)
}
