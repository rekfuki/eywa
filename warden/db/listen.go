package db

import (
	"context"

	"github.com/jackc/pgx"
)

// Subscription represents postgres notification subscription
type Subscription struct {
	client  *Client
	channel string
	conn    *pgx.Conn
	quit    chan struct{}
	ErrChan chan error
	Notify  chan *pgx.Notification
}

// Close closes subscription
func (s *Subscription) Close() {
	s.quit <- struct{}{}
	s.client.pool.Release(s.conn)
}

// Listen listens to notifications from postgres
func (c *Client) Listen(channel string) (*Subscription, error) {
	conn, err := c.pool.Acquire()
	if err != nil {
		return nil, err
	}

	_, err = conn.Exec("listen " + channel)
	if err != nil {
		return nil, err
	}

	errChan := make(chan error)
	notify := make(chan *pgx.Notification)
	quit := make(chan struct{})

	subscription := &Subscription{
		client:  c,
		channel: channel,
		conn:    conn,
		quit:    quit,
		ErrChan: errChan,
		Notify:  notify,
	}

	go func() {
		for {
			select {
			case <-quit:
				return
			default:
				for {
					notification, err := conn.WaitForNotification(context.Background())
					if err != nil {
						errChan <- err
					}
					notify <- notification
				}
			}
		}
	}()

	return subscription, nil
}
