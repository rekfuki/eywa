package broker

import (
	"net/http"
	"time"
)

const (
	// TimelineType represents timeline log type
	TimelineType = "timeline"

	// RequestType represents request log type
	RequestType = "request"

	// ExecutionType represents execution log type
	ExecutionType = "execution"

	// SyncExecType represents sync executuin log type
	SyncExecType = "sync-execution"

	// EventLogType represents event log type
	EventLogType = "event"
)

// Message represents message sent over nats
type Message struct {
	UserID    string    `json:"user_id"`
	RequestID string    `json:"request_id"`
	Timestamp time.Time `json:"timestamp"`
}

// MessageInterface ...
type MessageInterface interface {
	SetDefaults()
}

// SetDefaults ...
func (m Message) SetDefaults() {
	if m.Timestamp == (time.Time{}) {
		m.Timestamp = time.Now()
	}
}

// Log represents log message
type Log struct {
	Message
	TimelineLog *TimelineLog `json:"timeline_log,omitempty"`
	EventLog    *EventLog    `json:"event_log,omitempty"`
}

// TimelineLog represents timeline log message payload
type TimelineLog struct {
	FunctionID string    `json:"function_id"`
	EventName  string    `json:"event_name"`
	EventType  string    `json:"event_type"`
	Response   int       `json:"response"`
	Method     string    `json:"method"`
	Duration   int64     `json:"duration"`
	CreatedAt  time.Time `json:"created_at"`
}

// EventLog represents any events that occur during execution
type EventLog struct {
	Type         string                 `json:"type"`
	IsError      bool                   `json:"is_error"`
	FunctionName string                 `json:"function_name"`
	FunctionID   string                 `json:"function_id"`
	Message      string                 `json:"message"`
	Payload      map[string]interface{} `json:"payload,omitempty"`
	CreatedAt    time.Time              `json:"generated_at"`
}

// QueueRequestMessage for sending async request message via stan
type QueueRequestMessage struct {
	Message
	Payload QueueRequest
}

// QueueRequest for asynchronous processing
type QueueRequest struct {
	UserID       string      `json:"user_id"`
	RequestID    string      `jons:"request_id"`
	Headers      http.Header `json:"headers"`
	Body         []byte      `json:"body"`
	Path         string      `json:"path"`
	QueryParams  string      `json:"query"`
	FunctionID   string      `json:"function_id"`
	FunctionName string      `json:"function_name"`
	CallbackURL  string      `json:"callback_url"`
	QueuedAt     time.Time   `json:"queued_at"`
}
