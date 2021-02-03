package trigger

import (
	"fmt"
	"os"
	"time"
)

// Entry represents a triger entry
type Entry struct {
	Trigger *Trigger
	Data    Fields
	Type    Type
	Message string
	Time    time.Time
}

// WithFields sends a trigger with fields
func (entry *Entry) WithFields(fields Fields) *Entry {
	data := make(Fields, len(entry.Data)+len(fields))
	for k, v := range entry.Data {
		data[k] = v
	}
	for k, v := range fields {
		data[k] = v
	}

	return &Entry{Trigger: entry.Trigger, Data: data, Time: entry.Time}
}

// Fire fires the entry
func (entry *Entry) Fire(t Type, args ...interface{}) {
	entry.fire(t, fmt.Sprint(args...))
}

// FireForEach trigers all the hooks for each arg seperately
func (entry *Entry) FireForEach(t Type, args ...interface{}) {
	for _, arg := range args {
		entry.fire(t, fmt.Sprint(arg))
	}
}

func (entry Entry) fire(t Type, msg string) {
	entry.Time = time.Now()
	entry.Type = t
	entry.Message = msg

	err := entry.Trigger.Hooks.Fire(entry.Type, &entry)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to fire hook: %v\n", err)
	}
}
