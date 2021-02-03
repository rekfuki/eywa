package hooks

import (
	"time"

	"eywa/go-libs/broker"
	"eywa/go-libs/trigger"
)

// EventHook represents processing function for event trigger
func EventHook(entry *trigger.Entry) broker.MessageInterface {
	elm := broker.Log{
		Message: broker.Message{
			Timestamp: entry.Time,
		},
		EventLog: &broker.EventLog{},
	}

	if val, exists := entry.Data["request_id"]; exists {
		elm.Message.RequestID = val.(string)
	}

	if val, exists := entry.Data["user_id"]; exists {
		elm.Message.UserID = val.(string)
	}

	if val, exists := entry.Data["type"]; exists {
		elm.EventLog.Type = val.(string)
	}

	if val, exists := entry.Data["created_at"]; exists {
		elm.EventLog.CreatedAt = val.(time.Time)
	} else {
		elm.EventLog.CreatedAt = entry.Time
	}

	if val, exists := entry.Data["function_name"]; exists {
		elm.EventLog.FunctionName = val.(string)
	}

	if val, exists := entry.Data["function_id"]; exists {
		elm.EventLog.FunctionID = val.(string)
	}

	if val, exists := entry.Data["is_error"]; exists {
		elm.EventLog.IsError = val.(bool)
	}

	elm.EventLog.Message = entry.Message

	return elm
}
