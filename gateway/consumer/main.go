package main

import (
	"flag"
	"os"
	"time"

	"github.com/fvbock/endless"
	"github.com/kelseyhightower/envconfig"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/nats-io/stan.go"
	log "github.com/sirupsen/logrus"

	"eywa/gateway/clients/k8s"
	"eywa/gateway/consumer/listener"
	"eywa/gateway/metrics"
	"eywa/gateway/types"
	"eywa/go-libs/broker"
)

// Config represents gateway-queue configuration
type Config struct {
	NatsURL       string `envconfig:"nats_url" default:"nats://nats.nats:4222"`
	StanClusterID string `default:"stan"`
	StanClientID  string `envconfig:"gateway_queue_stan_client_id" default:"gateway-queue"`
	MaxInflight   int    `envconfig:"max_inflight" default:"100"`
	RetryCount    int    `envconfig:"retry_count" default:"3"`
	RetrySleep    int    `envconfig:"retry_sleep" default:"3"`
}

func main() {
	var conf Config
	err := envconfig.Process("", &conf)
	if err != nil {
		log.Fatalf("Failed to parse env: %s", err)
	}

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

	hostname, _ := os.Hostname()
	clientID := conf.StanClientID + broker.GetClientID(hostname)
	bc, err := broker.Connect(conf.NatsURL, conf.StanClusterID, clientID, 100, 5)
	if err != nil {
		log.Fatalf("Failed to setup nats-streaming broker: %s", err)
	}

	metrics := metrics.Setup(k8s, time.Second*5)

	e := echo.New()
	e.Use(middleware.Recover())

	// Expose metrics for prometheus
	e.GET("/metrics", echo.WrapHandler(metrics.PrometheusHandler()))

	listener := listener.New(&listener.Config{
		K8s:         k8s,
		Metrics:     metrics,
		Broker:      bc,
		MaxInFlight: conf.MaxInflight,
		RetryCount:  conf.RetryCount,
		RetrySleep:  conf.RetrySleep,
	})

	qSub, err := bc.QueueSubscribe(
		types.AsyncExecSubject, "gateway-consumer",
		listener.HandleMessage,
		// stan.MaxInflight(conf.MaxInflight), // TODO: Not sure if needed
		stan.DeliverAllAvailable(),
		stan.SetManualAckMode(),
		stan.DurableName("durable"))
	if err != nil {
		log.Fatalf("Failed to subscribe to topic %s: %s", types.AsyncExecSubject, err)
	}

	endless.DefaultHammerTime = 10 * time.Second
	endless.DefaultReadTimeOut = 295 * time.Second
	if err := endless.ListenAndServe(":8888", e); err != nil {
		log.Infof("Server stopped: %s", err)
	}

	qSub.Close()
}
