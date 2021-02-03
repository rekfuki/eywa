package broker

import (
	"eywa/go-libs/trigger"
)

// Processor represents processor function type
type Processor func(*trigger.Entry) MessageInterface

// LogHandler represents log handler object
type LogHandler struct {
	topic     string
	processor Processor
	broker    *Client
	async     bool
}

// NewLogHandler returns new log handler that produces to stan
func NewLogHandler(topic string, broker *Client, processor Processor, async bool) *LogHandler {
	return &LogHandler{
		topic:     topic,
		processor: processor,
		broker:    broker,
		async:     async,
	}
}

// Fire sends the log message to kafka
func (l *LogHandler) Fire(entry *trigger.Entry) error {
	msg := l.processor(entry)
	if msg == nil {
		return nil
	}

	if l.async {
		return l.broker.ProduceAsync(l.topic, msg)
	}
	return l.broker.ProduceSync(l.topic, msg)
}
