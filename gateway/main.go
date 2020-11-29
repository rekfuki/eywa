package main

import (
	"flag"
	"time"

	log "github.com/sirupsen/logrus"

	"eywa/gateway/clients/k8s"
	"eywa/gateway/metrics"
	"eywa/gateway/server"
)

func main() {
	inCluster := flag.Bool("in-cluster", true, "(optional) running inside the cluser")
	debug := flag.Bool("debug", false, "(optional) set log level to debug")
	flag.Parse()

	if *debug {
		log.SetLevel(log.DebugLevel)
	}

	k8s, err := k8s.Setup(&k8s.Config{
		InCluster:           *inCluster,
		CacheExpiryDuration: time.Second * 5,
	})
	if err != nil {
		log.Fatalf("Failed to setup k8s client: %s", err)
	}

	metrics := metrics.Setup(k8s, time.Second*5)
	go metrics.FunctionWatcher()

	params := &server.ContextParams{
		K8s:     k8s,
		Metrics: metrics,
	}

	server.Run(params)

	log.Exit(0)
}
