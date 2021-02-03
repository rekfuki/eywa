package hooks

import (
	"time"

	"eywa/go-libs/broker"
	"eywa/go-libs/trigger"
)

// TimelineHook represents processing function for timeline trigger
func TimelineHook(entry *trigger.Entry) broker.MessageInterface {
	tlm := broker.Log{
		Message: broker.Message{
			Timestamp: entry.Time,
		},
		TimelineLog: &broker.TimelineLog{},
	}

	if val, exists := entry.Data["request_id"]; exists {
		tlm.Message.RequestID = val.(string)
	}

	if val, exists := entry.Data["user_id"]; exists {
		tlm.Message.UserID = val.(string)
	}

	if val, exists := entry.Data["function_id"]; exists {
		tlm.TimelineLog.FunctionID = val.(string)
	}

	if val, exists := entry.Data["event_name"]; exists {
		tlm.TimelineLog.EventName = val.(string)
	}

	if val, exists := entry.Data["event_type"]; exists {
		tlm.TimelineLog.EventType = val.(string)
	}

	if val, exists := entry.Data["created_at"]; exists {
		tlm.TimelineLog.CreatedAt = val.(time.Time)
	} else {
		tlm.TimelineLog.CreatedAt = entry.Time
	}

	if val, exists := entry.Data["response"]; exists {
		tlm.TimelineLog.Response = val.(int)
	}

	if val, exists := entry.Data["method"]; exists {
		tlm.TimelineLog.Method = val.(string)
	}

	if val, exists := entry.Data["duration"]; exists {
		tlm.TimelineLog.Duration = val.(int64)
	}

	if val, exists := entry.Data["created_at"]; exists {
		tlm.TimelineLog.CreatedAt = val.(time.Time)
	} else {
		tlm.TimelineLog.CreatedAt = entry.Time
	}

	return tlm
}
