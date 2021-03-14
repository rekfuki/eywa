package metrics

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
	"gopkg.in/resty.v1"

	"eywa/gateway/clients/k8s"
)

// Client metrics client
type Client struct {
	metrics       *metrics
	services      []k8s.FunctionStatus
	watchInterval time.Duration
	k8sClient     *k8s.Client
	promrc        *resty.Client
}

type metrics struct {
	functionsHistogram        *prometheus.HistogramVec
	queueHistogram            *prometheus.HistogramVec
	functionInvocation        *prometheus.CounterVec
	functionInvocationStarted *prometheus.CounterVec
	serviceReplicasGauge      *prometheus.GaugeVec
}

// Setup sets up prometheus counters and histograms
func Setup(k8sClient *k8s.Client, prometheusURL *string, watchInterval time.Duration) *Client {
	client := &Client{
		metrics:       setupMetrics(),
		services:      []k8s.FunctionStatus{},
		watchInterval: watchInterval,
		k8sClient:     k8sClient,
	}

	if prometheusURL != nil {
		client.promrc = resty.New().
			SetHostURL(*prometheusURL + "/api/v1").
			SetLogger(ioutil.Discard).
			SetRetryCount(3).
			SetTimeout(10 * time.Second)
	}

	prometheus.MustRegister(client)
	return client
}

func setupMetrics() *metrics {
	gatewayFunctionsHistogram := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name: "gateway_function_duration_milliseconds",
		Help: "Function time taken",
	}, []string{"function_id", "function_name", "user_id"})

	gatewayAsyncQueueHistogram := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name: "gateway_queue_dwell_duration_milliseconds",
		Help: "Function time taken",
	}, []string{"function_id", "function_name", "user_id"})

	gatewayFunctionInvocation := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "gateway",
			Subsystem: "function",
			Name:      "invocation_total",
			Help:      "Function metrics",
		},
		[]string{"function_id", "function_name", "user_id", "path", "code"},
	)

	serviceReplicas := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "gateway",
			Name:      "service_count",
			Help:      "Service replicas",
		},
		[]string{"function_name"},
	)

	gatewayFunctionInvocationStarted := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "gateway",
			Subsystem: "function",
			Name:      "invocation_started",
			Help:      "The total number of function HTTP requests started.",
		},
		[]string{"function_id", "function_name", "user_id", "path", "method"},
	)

	return &metrics{
		functionsHistogram:        gatewayFunctionsHistogram,
		queueHistogram:            gatewayAsyncQueueHistogram,
		functionInvocation:        gatewayFunctionInvocation,
		functionInvocationStarted: gatewayFunctionInvocationStarted,
		serviceReplicasGauge:      serviceReplicas,
	}
}

// ObserveInvocationStarted records function invocation started metrics in Prometheus
func (c *Client) ObserveInvocationStarted(fnID, fnName, userID, path, method string) {
	c.metrics.functionInvocationStarted.With(prometheus.Labels{
		"function_id":   fnID,
		"function_name": fnName,
		"user_id":       userID,
		"path":          path,
		"method":        method,
	}).Inc()
}

// ObserveInvocationComplete records function invocation complete metrics in Prometheus
func (c *Client) ObserveInvocationComplete(fnID, fnName, userID, path string, statusCode int, duration time.Duration) {
	milliseconds := duration.Milliseconds()
	c.metrics.functionsHistogram.
		With(prometheus.Labels{
			"function_id":   fnID,
			"function_name": fnName,
			"user_id":       userID,
		}).
		Observe(float64(milliseconds))

	code := strconv.Itoa(statusCode)
	c.metrics.functionInvocation.
		With(prometheus.Labels{
			"function_id":   fnID,
			"function_name": fnName,
			"user_id":       userID,
			"path":          path,
			"code":          code,
		}).
		Inc()
}

// ObserveDwellTime records function dwell time in the queue metrics in Prometheus
func (c *Client) ObserveDwellTime(fnID, fnName, userID string, duration time.Duration) {
	milliseconds := duration.Milliseconds()
	c.metrics.queueHistogram.
		With(prometheus.Labels{
			"function_id":   fnID,
			"function_name": fnName,
			"user_id":       userID,
		}).
		Observe(float64(milliseconds))
}

// FunctionWatcher watches currently deployed functions and stores them for metrics
func (c *Client) FunctionWatcher() {
	for {
		functions, err := c.k8sClient.GetFunctionsStatus()
		if err != nil {
			log.Errorf("Failed to list current functions: %s", err)
			continue
		}

		c.services = functions

		time.Sleep(c.watchInterval)
	}
}

// Collect defines metrics collection function for prometheus
func (c *Client) Collect(ch chan<- prometheus.Metric) {
	c.metrics.functionInvocation.Collect(ch)
	c.metrics.functionsHistogram.Collect(ch)
	c.metrics.queueHistogram.Collect(ch)
	c.metrics.functionInvocationStarted.Collect(ch)
	c.metrics.serviceReplicasGauge.Reset()
	for _, service := range c.services {
		var serviceName string
		if len(service.Namespace) > 0 {
			serviceName = fmt.Sprintf("%s.%s", service.Name, service.Namespace)
		} else {
			serviceName = service.Name
		}
		c.metrics.serviceReplicasGauge.
			WithLabelValues(serviceName).
			Set(float64(service.Replicas))
	}
	c.metrics.serviceReplicasGauge.Collect(ch)
}

// Describe defines metrics description function for prometheus
func (c *Client) Describe(ch chan<- *prometheus.Desc) {
	c.metrics.functionInvocation.Describe(ch)
	c.metrics.queueHistogram.Describe(ch)
	c.metrics.functionsHistogram.Describe(ch)
	c.metrics.serviceReplicasGauge.Describe(ch)
	c.metrics.functionInvocationStarted.Describe(ch)
}

// PrometheusHandler returns prometheus handler
func (c *Client) PrometheusHandler() http.Handler {
	return promhttp.Handler()
}
