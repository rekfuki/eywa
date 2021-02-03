package broker

import (
	"encoding/json"
	"regexp"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/stan.go"
	log "github.com/sirupsen/logrus"
)

var supportedCharacters = regexp.MustCompile("[^a-zA-Z0-9-_]+")

// Client represents broker client wrapper around stan
type Client struct {
	natsURL        string
	clusterID      string
	clientID       string
	maxReconnect   int
	reconnectDelay time.Duration
	stan.Conn
}

// Connect connects to nats streaming
func Connect(natsURL, clusterID, clientID string, maxReconnect int, rcDelay time.Duration) (*Client, error) {
	log.Printf("Connect: %s\n", natsURL)

	broker := &Client{
		natsURL:        natsURL,
		clusterID:      clusterID,
		clientID:       clientID,
		maxReconnect:   maxReconnect,
		reconnectDelay: rcDelay,
	}

	nc, err := nats.Connect(natsURL, nats.MaxReconnects(-1))
	if err != nil {
		log.Fatalf("Failed to connect to nats: %s", err)
	}

	sc, err := stan.Connect(
		clusterID,
		clientID,
		stan.NatsConn(nc),
		stan.SetConnectionLostHandler(func(_ stan.Conn, err error) {
			log.Errorf("Disconnected from %s", natsURL)
			broker.reconnect()
		}))
	if err != nil {
		return nil, err
	}

	log.Infof("Connected to: %s", natsURL)
	broker.Conn = sc

	return broker, err
}

func (c *Client) reconnect() {
	log.Printf("Reconnect\n")

	for i := 0; i < c.maxReconnect; i++ {
		newClient, err := Connect(c.natsURL, c.clusterID, c.clientID, c.maxReconnect, c.reconnectDelay)
		if err == nil {
			c = newClient
			log.Printf("Reconnecting (%d/%d) to %s succeeded\n", i+1, c.maxReconnect, c.natsURL)

			return
		}

		nextTryIn := (time.Duration(i+1) * c.reconnectDelay).String()

		log.Printf("Reconnecting (%d/%d) to %s failed\n", i+1, c.maxReconnect, c.natsURL)
		log.Printf("Waiting %s before next try", nextTryIn)

		time.Sleep(time.Duration(i) * c.reconnectDelay)
	}

	log.Printf("Reconnecting limit (%d) reached\n", c.maxReconnect)
}

// ProduceSync produces sync message on a given topic
func (c *Client) ProduceSync(topic string, msg MessageInterface) error {
	msg.SetDefaults()
	bytes, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	return c.Publish(topic, bytes)
}

// ProduceAsync produces assync message on a given topic
func (c *Client) ProduceAsync(topic string, msg MessageInterface) error {
	msg.SetDefaults()
	bytes, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	_, err = c.PublishAsync(topic, bytes, func(ackedNuid string, err error) {
		if err != nil {
			log.Warnf("Warning: error publishing msg id %s: %v\n", ackedNuid, err.Error())
		}
	})
	if err != nil {
		return err
	}

	return nil
}

// Close closes the queue
func (c *Client) Close() error {
	return c.Conn.Close()
}

// GetClientID returns sanitsed client ID
func GetClientID(value string) string {
	return supportedCharacters.ReplaceAllString(value, "_")
}
