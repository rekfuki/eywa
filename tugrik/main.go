package main

import (
	"github.com/kelseyhightower/envconfig"
	log "github.com/sirupsen/logrus"

	"eywa/tugrik/clients/gateway"
	"eywa/tugrik/db"
	"eywa/tugrik/server"
)

// Config represents wardends startup config
type Config struct {
	DB            db.Config
	GatewayAPIURL string `envconfig:"gateway_api_url" default:"http://gateway-api.faas-system:8080"`
}

func main() {
	var conf Config
	err := envconfig.Process("", &conf)
	if err != nil {
		log.Fatalf("Failed to parse env: %s", err)
	}

	db, err := db.Connect(conf.DB)
	if err != nil {
		log.Fatalf("Failed to connect to database: %s", err)
	}

	gateway := gateway.New(conf.GatewayAPIURL)

	params := &server.ContextParams{
		DB:      db,
		Gateway: gateway,
	}

	server.Run(params)
}
