package listener

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/nats-io/stan.go"
	log "github.com/sirupsen/logrus"
	"gopkg.in/resty.v1"

	ett "eywa/execution-tracker/types"
	"eywa/gateway/clients/k8s"
	"eywa/gateway/hooks"
	"eywa/gateway/metrics"
	"eywa/gateway/types"
	"eywa/go-libs/broker"
	"eywa/go-libs/trigger"
	wet "eywa/watchdog/executor"
)

// Config listener configuration
type Config struct {
	K8s         *k8s.Client
	Metrics     *metrics.Client
	Broker      *broker.Client
	MaxInFlight int
	RetryCount  int
	RetrySleep  int
}

// Listener listens on history topic and inserts records into mongo db
type Listener struct {
	k8s        *k8s.Client
	metrics    *metrics.Client
	broker     *broker.Client
	incMsg     chan *stan.Msg
	rc         *resty.Client
	retryCount int
	retrySleep int
}

// New Listener
func New(conf *Config) *Listener {
	rc := resty.New()
	rc.SetTransport(&http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 0,
		}).DialContext,

		MaxIdleConns:          1,
		DisableKeepAlives:     true,
		IdleConnTimeout:       120 * time.Millisecond,
		ExpectContinueTimeout: 1500 * time.Millisecond,
	})

	listener := &Listener{
		k8s:        conf.K8s,
		metrics:    conf.Metrics,
		broker:     conf.Broker,
		incMsg:     make(chan *stan.Msg),
		rc:         rc,
		retryCount: conf.RetryCount,
		retrySleep: conf.RetrySleep,
	}

	eventLogHandler := broker.NewLogHandler(types.LogsSubject, conf.Broker, hooks.EventHook, false)
	timelineLogHandler := broker.NewLogHandler(types.LogsSubject, conf.Broker, hooks.TimelineHook, false)

	trigger.AddHook(eventLogHandler, []trigger.Type{types.EventHookType})
	trigger.AddHook(timelineLogHandler, []trigger.Type{types.TimelineHookType})

	for i := 0; i < conf.MaxInFlight; i++ {
		go func() {
			for msg := range listener.incMsg {
				listener.process(msg)
			}
		}()
	}

	return listener
}

// HandleMessage handle messages from STAN
func (l *Listener) HandleMessage(msg *stan.Msg) {
	l.incMsg <- msg
}

// Process read from message channel and handle the messages
func (l *Listener) process(msg *stan.Msg) {
	if err := msg.Ack(); err != nil {
		log.Errorf("Failed to ack message %s: %s", msg.String(), err)
	}

	qrm := broker.QueueRequestMessage{}
	if err := json.Unmarshal(msg.Data, &qrm); err != nil {
		log.Errorf("Failed to unmarshal queue request. Error: %s. Data: ", err, string(msg.Data))
		return
	}

	req := qrm.Payload

	if !strings.HasPrefix(req.Path, "/") {
		req.Path = "/" + req.Path
	}

	if err := validateMessage(req); err != nil {
		log.Errorf("%s. Dropping...", err.Error())
		return
	}

	defaultTimelineFields := trigger.Fields{
		"user_id":     req.UserID,
		"request_id":  req.RequestID,
		"function_id": req.FunctionID,
	}

	defaultEventFields := trigger.Fields{
		"user_id":       req.UserID,
		"request_id":    req.RequestID,
		"type":          ett.EventTypeSystem,
		"function_name": req.FunctionName,
		"function_id":   req.FunctionID,
	}

	trigger.WithFields(defaultTimelineFields).WithFields(trigger.Fields{
		"event_name": "Dwell Time",
		"event_type": ett.TimelineEventTypeDequeued,
		"response":   http.StatusOK,
		"duration":   time.Since(req.QueuedAt).Milliseconds(),
	}).Fire(types.TimelineHookType)

	l.metrics.ObserveDwellTime(req.FunctionID, req.FunctionName, req.UserID, time.Since(req.QueuedAt))

	sleepDuration := time.Minute * 0
	for attempt := 1; attempt <= 3; attempt++ {
		time.Sleep(sleepDuration)
		sleepDuration = sleepDuration + time.Minute*3

		started := time.Now()

		trigger.WithFields(defaultEventFields).WithFields(trigger.Fields{
			"path":         req.Path,
			"query_params": req.QueryParams,
			"body":         req.Body,
			"headers":      req.Headers,
			"message":      types.AsyncExecutionStartMessage(attempt, req.FunctionName),
		}).Fire(types.EventHookType)

		l.metrics.ObserveInvocationStarted(req.FunctionID, req.FunctionName,
			req.UserID, req.Path, http.MethodPost)

		filter := k8s.LabelSelector().
			Equals(types.UserIDLabel, req.UserID).
			Equals(types.FunctionIDLabel, req.FunctionID)
		scaleResult, err := l.k8s.ScaleFromZero(filter)
		if err != nil {
			log.Errorf("Failed to scale function %q from zero: %s", req.FunctionID, err)

			trigger.WithFields(defaultTimelineFields).WithFields(trigger.Fields{
				"event_name": fmt.Sprintf("Attempt #%d", attempt),
				"event_type": ett.TimelineEventTypeSystemError,
				"response":   http.StatusServiceUnavailable,
				"duration":   time.Since(started).Milliseconds(),
			}).Fire(types.TimelineHookType)

			trigger.WithFields(defaultEventFields).
				WithFields(trigger.Fields{
					"is_error": true,
					"message":  types.ServerErrorMessage(),
				}).Fire(types.EventHookType)

			continue
		}

		if !scaleResult.Found {
			log.Errorf("Function %q deployment not found")

			trigger.WithFields(defaultTimelineFields).WithFields(trigger.Fields{
				"event_name": fmt.Sprintf("Attempt #%d", attempt),
				"event_type": ett.TimelineEventTypeFailed,
				"response":   http.StatusNotFound,
				"duration":   time.Since(started).Milliseconds(),
			}).Fire(types.TimelineHookType)

			trigger.WithFields(defaultEventFields).
				WithFields(trigger.Fields{
					"is_error": true,
					"message":  types.FunctionNotFoundMessage(req.RequestID, req.FunctionID, req.FunctionName),
				}).Fire(types.EventHookType)

			continue
		}

		if !scaleResult.Available {
			log.Errorf("Function %q scale request timed-out after %fs", req.FunctionID, scaleResult.Duration)

			trigger.WithFields(defaultTimelineFields).WithFields(trigger.Fields{
				"event_name": fmt.Sprintf("Attempt #%d", attempt),
				"event_type": ett.TimelineEventTypeSystemError,
				"response":   http.StatusServiceUnavailable,
				"duration":   time.Since(started).Milliseconds(),
			}).Fire(types.TimelineHookType)

			trigger.WithFields(defaultEventFields).
				WithFields(trigger.Fields{
					"is_error": true,
					"message":  types.ServerErrorMessage(),
				}).Fire(types.EventHookType)

			continue
		}

		start := time.Now()
		functionAddr, err := l.k8s.Resolve(req.FunctionID)
		if err != nil {
			log.Errorf("k8s error: cannot find %s: %s\n", req.FunctionID, err)

			trigger.WithFields(defaultTimelineFields).WithFields(trigger.Fields{
				"event_name": fmt.Sprintf("Attempt #%d", attempt),
				"event_type": ett.TimelineEventTypeSystemError,
				"response":   http.StatusServiceUnavailable,
				"duration":   time.Since(started).Milliseconds(),
			}).Fire(types.TimelineHookType)

			trigger.WithFields(defaultEventFields).
				WithFields(trigger.Fields{
					"is_error": true,
					"message":  types.ServerErrorMessage(),
				}).Fire(types.EventHookType)

			continue
		}

		trigger.WithFields(defaultTimelineFields).WithFields(trigger.Fields{
			"event_name": fmt.Sprintf("Attempt #%d", attempt),
			"event_type": ett.TimelineEventTypeRunning,
			"response":   http.StatusOK,
			"duration":   int64(0),
		}).Fire(types.TimelineHookType)

		headers := map[string]string{}
		for k, h := range req.Headers {
			headers[k] = strings.Join(h, ",")
		}

		url := functionAddr + req.Path
		var result wet.FunctionResponse
		functionRes, err := l.rc.R().
			SetBody(req.Body).
			SetResult(&result).
			SetHeaders(headers).
			SetQueryString(req.QueryParams).
			Post(url)
		if err != nil || functionRes.IsError() {
			log.Errorf("Failed to execute function request [%s] %q: %s", http.MethodPost, url, err)

			trigger.WithFields(defaultTimelineFields).WithFields(trigger.Fields{
				"event_name": fmt.Sprintf("Attempt #%d", attempt),
				"event_type": ett.TimelineEventTypeSystemError,
				"response":   http.StatusServiceUnavailable,
				"duration":   time.Since(started).Milliseconds(),
			}).Fire(types.TimelineHookType)

			trigger.WithFields(defaultEventFields).
				WithFields(trigger.Fields{
					"is_error": true,
					"message":  types.ServerErrorMessage(),
				}).Fire(types.EventHookType)

			continue
		}

		duration := time.Since(start)
		l.metrics.ObserveInvocationComplete(req.FunctionID, req.FunctionName, req.UserID, req.Path, result.Status, duration)

		log.Infof("[Attempt: #%s] Invoked: %s-%s [%d] in %fs", attempt, req.FunctionID,
			req.FunctionName, result.Status, duration.Seconds())

		eventType := ett.TimelineEventTypeFinished
		if result.Status >= 400 {
			eventType = ett.TimelineEventTypeFailed
		}

		trigger.WithFields(defaultTimelineFields).WithFields(trigger.Fields{
			"event_name": fmt.Sprintf("Attempt #%d", attempt),
			"event_type": eventType,
			"response":   result.Status,
			"duration":   duration.Milliseconds(),
		}).Fire(types.TimelineHookType)

		defaultEventFields["type"] = ett.EventTypeUser
		defaultEventFields["created_at"] = started.Add(duration)

		// Inject back request id header
		if result.Headers == nil {
			result.Headers = http.Header{}
		}

		result.Headers.Set("X-Request-Id", req.RequestID)
		result.Headers.Del("X-Eywa-Token")

		trigger.WithFields(defaultEventFields).WithFields(trigger.Fields{
			"is_error": eventType == ett.TimelineEventTypeFailed,
			"status":   result.Status,
			"headers":  result.Headers,
			"body":     result.Body,
			"stdout":   result.Stdout,
			"stderr":   result.Stderr,
			"message":  types.AsyncExecutionFinishMessage(req.FunctionName, attempt, result.Status, duration),
		}).Fire(types.EventHookType)

		if req.CallbackURL != "" {
			log.Infof("Sending callback to: %s\n", req.CallbackURL)
			_, err := l.rc.R().
				SetHeaders(map[string]string{
					"X-Function-Name":   req.FunctionName,
					"X-Function-Id":     req.FunctionID,
					"X-Function-Status": fmt.Sprint(result.Status),
				}).
				SetBody(functionRes.Body()).
				Post(req.CallbackURL)
			if err != nil {
				log.Warnf("Failed call callback url %q: %s", req.CallbackURL, err)
				trigger.WithFields(defaultEventFields).WithFields(trigger.Fields{
					"is_error": true,
					"message":  types.CallbackError(err.Error()),
				}).Fire(types.EventHookType)
			}
		}

		if eventType != ett.TimelineEventTypeFinished {
			continue
		}

		return
	}
}

func validateMessage(req broker.QueueRequest) error {
	if req.UserID == "" {
		return fmt.Errorf("message is missing User ID")
	}

	if req.RequestID == "" {
		return fmt.Errorf("message is missing Request ID")
	}

	if req.FunctionID == "" {
		return fmt.Errorf("message is missing Function ID")
	}

	if req.FunctionName == "" {
		return fmt.Errorf("message is missing Function Name")
	}

	return nil
}
