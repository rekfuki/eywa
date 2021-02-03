package metrics

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"

	"eywa/gateway/clients/k8s"
)

// Client metrics client
type Client struct {
	metrics       *metrics
	services      []k8s.FunctionStatus
	watchInterval time.Duration
	k8sClient     *k8s.Client
}

type metrics struct {
	functionsHistogram        *prometheus.HistogramVec
	functionInvocation        *prometheus.CounterVec
	functionInvocationStarted *prometheus.CounterVec
	serviceReplicasGauge      *prometheus.GaugeVec
}

// Setup sets up prometheus counters and histograms
func Setup(k8sClient *k8s.Client, watchInterval time.Duration) *Client {
	gatewayFunctionsHistogram := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name: "gateway_functions_seconds",
		Help: "Function time taken",
	}, []string{"function_name", "user_id"})

	gatewayFunctionInvocation := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "gateway",
			Subsystem: "function",
			Name:      "invocation_total",
			Help:      "Function metrics",
		},
		[]string{"function_name", "user_id", "code"},
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
		[]string{"function_name"},
	)

	client := &Client{
		metrics: &metrics{
			functionsHistogram:        gatewayFunctionsHistogram,
			functionInvocation:        gatewayFunctionInvocation,
			functionInvocationStarted: gatewayFunctionInvocationStarted,
			serviceReplicasGauge:      serviceReplicas,
		},
		services:      []k8s.FunctionStatus{},
		watchInterval: watchInterval,
		k8sClient:     k8sClient,
	}

	prometheus.MustRegister(client)
	return client
}

// Observe records metrics in Prometheus
func (c *Client) Observe(method, fnName, userID string, statusCode int, event string, duration time.Duration) {
	switch event {
	case "completed":
		seconds := duration.Seconds()
		c.metrics.functionsHistogram.
			With(prometheus.Labels{"function_name": fnName, "user_id": userID}).
			Observe(seconds)

		code := strconv.Itoa(statusCode)

		c.metrics.functionInvocation.
			With(prometheus.Labels{"function_name": fnName, "user_id": userID, "code": code}).
			Inc()
	case "started":
		c.metrics.functionInvocationStarted.WithLabelValues(fnName).Inc()
	default:
		log.Errorf("Unknown metrics event %q", event)
	}
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
	c.metrics.functionsHistogram.Describe(ch)
	c.metrics.serviceReplicasGauge.Describe(ch)
	c.metrics.functionInvocationStarted.Describe(ch)
}

// PrometheusHandler returns prometheus handler
func (c *Client) PrometheusHandler() http.Handler {
	return promhttp.Handler()
}
