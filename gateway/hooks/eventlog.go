package hooks

import (
	"fmt"
	"time"

	"eywa/go-libs/broker"
	"eywa/go-libs/trigger"
)

const maxFieldSize = 4 * 2000 // 4 bytes of 2000 chars is the limit

// EventHook represents processing function for event trigger
func EventHook(entry *trigger.Entry) broker.MessageInterface {
	elm := broker.Log{
		Message: broker.Message{
			Timestamp: entry.Time,
		},
		EventLog: &broker.EventLog{
			Payload: make(map[string]interface{}),
		},
	}

	for k, v := range entry.Data {
		elm.EventLog.CreatedAt = entry.Time
		switch k {
		case "request_id":
			elm.Message.RequestID = v.(string)
		case "user_id":
			elm.Message.UserID = v.(string)
		case "type":
			elm.EventLog.Type = v.(string)
		case "created_at":
			elm.EventLog.CreatedAt = v.(time.Time)
		case "function_name":
			elm.EventLog.FunctionName = v.(string)
		case "function_id":
			elm.EventLog.FunctionID = v.(string)
		case "is_error":
			elm.EventLog.IsError = v.(bool)
		case "message":
			elm.EventLog.Message = v.(string)
		case "body", "stdout", "stderr":
			val, ok := v.([]uint8)
			if ok {
				if len(val) > maxFieldSize {
					v = fmt.Sprintf("Value too large to display (%d bytes), max allowed %d", len(val), maxFieldSize)
				}
				elm.EventLog.Payload[k] = v
			}
		default:
			elm.EventLog.Payload[k] = v
		}
	}

	return elm
}
