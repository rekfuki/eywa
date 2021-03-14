package types

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

// Timeline event types
const (
	TimelineEventTypeCreated        = "created"
	TimelineEventTypeQueued         = "queued"
	TimelineEventTypeDequeued       = "dequeued"
	TimelineEventTypeRunning        = "running"
	TimelineEventTypeFinished       = "finished"
	TimelineEventTypeFailed         = "failed"
	TimelineEventTypeCallbackFailed = "callback-failed"
	TimelineEventTypeSystemError    = "system-error"

	EventTypeSystem = "system"
	EventTypeUser   = "user"
)

// Allowed types for fast validation in consumer
var (
	AllowedEventTypes = map[string]struct{}{
		EventTypeSystem: {},
		EventTypeUser:   {},
	}

	AllowedTimelineEvents = map[string]struct{}{
		TimelineEventTypeCreated:     {},
		TimelineEventTypeQueued:      {},
		TimelineEventTypeRunning:     {},
		TimelineEventTypeFinished:    {},
		TimelineEventTypeFailed:      {},
		TimelineEventTypeSystemError: {},
		TimelineEventTypeDequeued:    {},
	}
)

// Headers type for pq marshaling
type Headers map[string]string

// Value returns marshaled headers
func (h Headers) Value() (driver.Value, error) {
	return json.Marshal(h)
}

// Scan decodes postgres value into a Headers type
func (h *Headers) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(b, &h)
}

// TimelineLog represents timeline log entry
type TimelineLog struct {
	RequestID  string    `json:"request_id" db:"request_id"`
	UserID     string    `json:"user_id" db:"user_id"`
	FunctionID string    `json:"function_id" db:"function_id"`
	EventType  string    `json:"event_type" db:"event_type"`
	EventName  string    `json:"event_name" db:"event_name"`
	Response   int       `json:"response" db:"response"`
	Method     string    `json:"method" db:"method"`
	Duration   int64     `json:"duration" db:"duration"`
	Timestamp  time.Time `json:"created_at" db:"timestamp"`
	ExpiresAt  time.Time `json:"-" db:"expires_at"`
}

// TimelineLogsResponse represents response when multiple timelines are returned
type TimelineLogsResponse struct {
	Page       int             `json:"page"`
	PerPage    int             `json:"per_page"`
	TotalCount int             `json:"total_count"`
	Objects    []TimelineBrief `json:"objects"`
}

// EventLogPayload custom type to be stored/scaned as jsonb
type EventLogPayload map[string]interface{}

// Value returns marshaled headers
func (elp EventLogPayload) Value() (driver.Value, error) {
	return json.Marshal(elp)
}

// Scan unmarshals jsonb in postgres to map[string]interface
func (elp *EventLogPayload) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(b, &elp)
}

// EventLog represents event log entry
type EventLog struct {
	ID           string          `json:"id" db:"id"`
	RequestID    string          `json:"request_id" db:"request_id"`
	UserID       string          `json:"user_id" db:"user_id"`
	Type         string          `json:"-" db:"type"`
	FunctionName string          `json:"function_name" db:"function_name"`
	FunctionID   string          `json:"function_id" db:"function_id"`
	Message      string          `json:"message" db:"message"`
	IsError      bool            `json:"is_error" db:"is_error"`
	Timestamp    time.Time       `json:"created_at" db:"timestamp"`
	ExpiresAt    time.Time       `json:"-" db:"expires_at"`
	Payload      EventLogPayload `json:"details" db:"payload"`
}

// EventLogsResponse represents response when multiple events are returned
type EventLogsResponse struct {
	Page       int        `json:"page"`
	PerPage    int        `json:"per_page"`
	TotalCount int        `json:"total_count"`
	Objects    []EventLog `json:"objects"`
}

// TimelineBrief is returned as small piece of information for timeline list
type TimelineBrief struct {
	RequestID    string    `json:"request_id"`
	FunctionID   string    `json:"function_id"`
	FunctionName string    `json:"function_name"`
	Status       int       `json:"status"`
	Age          time.Time `json:"age"`
	Duration     int64     `json:"duration"`
	IsError      bool      `json:"is_error"`
}

// TimelineDetails returns details about the timeline
type TimelineDetails struct {
	RequestID string         `json:"request_id"`
	Method    string         `json:"method"`
	Response  int            `json:"response"`
	Duration  int64          `json:"duration"`
	Age       time.Time      `json:"age"`
	Events    []EventDetails `json:"events"`
}

// EventDetails represents an event that occurs durint the timeline
type EventDetails struct {
	Name      string    `json:"name"`
	Response  int       `json:"response"`
	Duration  int64     `json:"duration"`
	IsError   bool      `json:"is_error"`
	Timestamp time.Time `json:"timestamp"`
}

// EventLogsQuery API request payload for filtering search
type EventLogsQuery struct {
	UserID       string    `json:"-"`
	FunctionID   string    `json:"function_id"`
	RequestID    string    `json:"request_id"`
	Query        string    `json:"query"`
	OnlyErrors   bool      `json:"only_errors"`
	Level        string    `json:"level" enum:"all,user,system"`
	TimestampMin time.Time `json:"timestamp_min,omitempty"`
	TimestampMax time.Time `json:"timestamp_max,omitempty"`
}
