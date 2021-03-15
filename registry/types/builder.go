package types

import (
	"time"

	"github.com/lib/pq"
)

// Build represents a build inside db
type Build struct {
	ImageID   string         `db:"image_id"`
	UserID    string         `db:"user_id"`
	Logs      pq.StringArray `db:"logs"`
	State     string         `db:"state"`
	CreatedAt time.Time      `db:"created_at"`
}
