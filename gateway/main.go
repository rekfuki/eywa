package main

import (
	"flag"

	log "github.com/sirupsen/logrus"

	"eywa/gateway/cache"
	"eywa/gateway/clients/k8s"
	"eywa/gateway/server"
)

func main() {
	inCluster := flag.Bool("in-cluster", true, "(optional) running inside the cluser")
	debug := flag.Bool("debug", false, "(optional) set log level to debug")
	flag.Parse()

	if *debug {
		log.SetLevel(log.DebugLevel)
	}

	cache := cache.Setup(&cache.Config{})

	k8s, err := k8s.Setup(&k8s.Config{
		InCluster: *inCluster,
	})
	if err != nil {
		log.Fatalf("Failed to setup k8s client: %s", err)
	}

	params := &server.ContextParams{
		Cache: cache,
		K8s:   k8s,
	}

	server.Run(params)

	log.Exit(0)
}
