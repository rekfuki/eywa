package main

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/kelseyhightower/envconfig"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/stan.go"
	uuid "github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"

	"eywa/execution-tracker/consumer/listener"
	"eywa/execution-tracker/db"
)

// Config is a populated by env variables and Vault
type Config struct {
	NatsURL           string        `envconfig:"nats_url" default:"nats://nats.nats:4222"`
	StanClusterID     string        `default:"stan"`
	StanClientID      string        `default:"execution-tracker-consumer"`
	BatchSize         int           `envconfig:"execution_tracker_batch_size" default:"3000"`
	FlushSeconds      int           `envconfig:"execution_tracker_flush_seconds" default:"2"`
	ExpireLogsMinutes time.Duration `envconfig:"execution_tracker_expire_logs" default:"10080m"`
	Postgres          db.Config
}

func main() {
	var conf Config
	err := envconfig.Process("", &conf)
	if err != nil {
		log.Fatalf("Failed to parse env: %s", err)
	}

	db, err := db.Connect(conf.Postgres)
	if err != nil {
		log.Fatalf("failed to connect to execution-tracker db: %s", err)
	}

	// delete expired logs
	go func() {
		t := time.NewTicker(time.Minute * time.Duration(conf.ExpireLogsMinutes))
		for range t.C {
			c, err := db.Begin()
			if err != nil {
				log.Fatalf("Failed to begin transaction: %s", err)
			}

			err = c.ExpireLogs()
			if err != nil {
				log.Fatalf("Failed to delete expired logs: %s", err)
			}

			err = c.Commit()
			if err != nil {
				log.Fatalf("Failed to commit transaction: %s", err)
			}

		}
	}()

	nc, err := nats.Connect(conf.NatsURL, nats.MaxReconnects(-1))
	if err != nil {
		log.Fatalf("Failed to connect to nats: %s", err)
	}

	id, _ := uuid.NewV4()
	clientID := conf.StanClientID + id.String()[:8]
	sc, err := stan.Connect(conf.StanClusterID, clientID, stan.NatsConn(nc))
	if err != nil {
		log.Fatalf("Failed to connect to STAN: %s", err)
	}

	listener := listener.New(db, conf.ExpireLogsMinutes, conf.BatchSize, conf.FlushSeconds)
	go listener.Process()

	qsub, err := sc.QueueSubscribe("logs", "execution-tracker-consumer", listener.HandleMessage, stan.DurableName("durable"))
	if err != nil {
		log.Fatalf("Failed to subscribe to logs topic: %s", err)
	}

	wait := make(chan os.Signal, 1)
	signal.Notify(wait, syscall.SIGINT, syscall.SIGTERM)
	<-wait

	log.Debug("Received Interrupt, shutting down")

	qsub.Close()
}
