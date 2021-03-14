package listener

import (
	"encoding/json"
	"time"

	"github.com/nats-io/stan.go"
	log "github.com/sirupsen/logrus"

	"eywa/execution-tracker/db"
	"eywa/execution-tracker/types"
	"eywa/go-libs/broker"
)

// Listener listens on history topic and inserts records into mongo db
type Listener struct {
	db              *db.Client
	batchSize       int
	flushSeconds    int
	expireIn        time.Duration
	incomingMessage chan *stan.Msg
}

// New Listener
func New(db *db.Client, expireIn time.Duration, batchSize, flushSeconds int) *Listener {

	listener := &Listener{
		db:              db,
		batchSize:       batchSize,
		flushSeconds:    flushSeconds,
		expireIn:        expireIn,
		incomingMessage: make(chan *stan.Msg, batchSize),
	}

	return listener
}

// HandleMessage handle messages from STAN
func (l *Listener) HandleMessage(msg *stan.Msg) {
	l.incomingMessage <- msg
}

// Process bacthed messages from nsq
func (l *Listener) Process() {
	batch := make([]*stan.Msg, 0, l.batchSize)
	flushTimer := time.NewTicker(time.Duration(l.flushSeconds) * time.Second)
	defer flushTimer.Stop()

	for {
		select {
		case msg := <-l.incomingMessage:
			batch = append(batch, msg)
			if len(batch) >= l.batchSize {
				l.batchInsert(batch)
				batch = batch[:0]
			}
		case <-flushTimer.C:
			if len(batch) > 0 {
				l.batchInsert(batch)
				batch = batch[:0]
			}
		}
	}
}

func (l *Listener) batchInsert(batch []*stan.Msg) {
	timelineLogs := []types.TimelineLog{}
	eventLogs := []types.EventLog{}

	for _, v := range batch {
		var msg broker.Log
		if err := json.Unmarshal(v.Data, &msg); err != nil {
			log.Warnf("Failed to unmarshal nsq message: %s", err)
			continue
		}

		expiresAt := msg.Timestamp.Add(l.expireIn)

		if msg.EventLog != nil {
			el := msg.EventLog
			if _, ok := types.AllowedEventTypes[el.Type]; !ok {
				log.Warnf("Unsupported event type %s", el.Type)
			}

			logEntry := types.EventLog{
				RequestID:    msg.RequestID,
				UserID:       msg.UserID,
				Type:         el.Type,
				FunctionName: el.FunctionName,
				FunctionID:   el.FunctionID,
				Message:      el.Message,
				Payload:      el.Payload,
				IsError:      el.IsError,
				Timestamp:    el.CreatedAt,
				ExpiresAt:    expiresAt,
			}
			eventLogs = append(eventLogs, logEntry)
		}

		if msg.TimelineLog != nil {
			tl := msg.TimelineLog
			if _, ok := types.AllowedTimelineEvents[tl.EventType]; !ok {
				log.Warnf("Unsupported timeline event type %s", tl.EventType)
			}

			if tl.Method == "" {
				tl.Method = "---"
			}

			logEntry := types.TimelineLog{
				RequestID:  msg.RequestID,
				UserID:     msg.UserID,
				FunctionID: tl.FunctionID,
				EventName:  tl.EventName,
				EventType:  tl.EventType,
				Response:   tl.Response,
				Method:     tl.Method,
				Duration:   tl.Duration,
				Timestamp:  tl.CreatedAt,
				ExpiresAt:  expiresAt,
			}
			timelineLogs = append(timelineLogs, logEntry)
		}
	}

	if len(timelineLogs) > 0 {
		total, dbErr := l.db.BulkInsertTimelineLogs(timelineLogs)
		if dbErr != nil {
			log.Errorf("Error bulk inserting timeline logs: %s", dbErr)
			return
		}
		log.Infof("Inserted %d timeline logs", total)
	}

	if len(eventLogs) > 0 {
		total, dbErr := l.db.BulkInsertEventLogs(eventLogs)
		if dbErr != nil {
			log.Errorf("Error bulk inserting event logs: %s", dbErr)
			return
		}
		log.Infof("Inserted %d event logs", total)
	}
}
