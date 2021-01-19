package main

import (
	"fmt"
	"strconv"
	"time"

	"github.com/kelseyhightower/envconfig"
	log "github.com/sirupsen/logrus"

	gwt "eywa/gateway/types"
	"eywa/idler/clients/gateway"
	"eywa/idler/clients/prometheus"
)

// Config represents gateway startup configuration
type Config struct {
	GatewayURL         string        `envconfig:"gateway_url" default:"http://gateway.faas-system:8080"`
	PrometheusURL      string        `envconfig:"prometheus_url" default:"http://linkerd-prometheus.linkerd:9090"`
	InactivityDuration time.Duration `envconfig:"inactivity_duration" default:"5m"`
}

// Idler represents the idler object
type Idler struct {
	gateway            *gateway.Client
	prometheus         *prometheus.Client
	reconcileInterval  time.Duration
	inactivityDuration time.Duration
}

func main() {
	var conf Config
	err := envconfig.Process("", &conf)
	if err != nil {
		log.Fatalf("Failed to parse env: %s", err)
	}

	idler := Idler{
		gateway:            gateway.New(conf.GatewayURL),
		prometheus:         prometheus.New(conf.PrometheusURL),
		reconcileInterval:  time.Second * 30,
		inactivityDuration: conf.InactivityDuration,
	}

	idler.Reconcile()
}

// Reconcile runs ilder recon loop
func (i *Idler) Reconcile() {
	for {
		functions, err := i.gateway.GetFunctions()
		if err != nil {
			log.Errorf("Failed to get functions: %s", err)
			continue
		}

		metrics := i.buildMetricsMap(functions)

		for _, fn := range functions {
			if v, found := metrics[fn.ID]; found {
				if v == float64(0) {
					log.Infof("%s\tidle\n", fn.Name)

					if fn.AvailableReplicas > 0 && fn.MinReplicas <= 0 {
						err := i.gateway.ScaleFunction(fn.ID, 0)
						if err != nil {
							log.Errorf("Failed to scale function: %s", err)
						}
					}

				} else {
					log.Infof("%s\tactive: %f\n", fn.Name, v)
				}
			}
		}
		time.Sleep(i.reconcileInterval)
	}
}

func (i *Idler) buildMetricsMap(functions []gwt.FunctionStatusResponse) map[string]float64 {
	metrics := make(map[string]float64)
	duration := fmt.Sprintf("%dm", int(i.inactivityDuration.Minutes()))
	for _, function := range functions {
		query := `sum(rate(gateway_function_invocation_total{function_name="` + function.ID + `", code=~".*"}[` + duration + `])) by (code, function_name)`
		res, err := i.prometheus.QueryMetrics(query)
		if err != nil {
			log.Errorf("Failed to get metrics from Prometheus: %s", err)
			continue
		}

		if len(res.Data.Result) > 0 {
			for _, v := range res.Data.Result {
				fmt.Println(v)
				if v.Metric.FunctionName == function.ID {
					metricValue := v.Value[1]
					switch metricValue.(type) {
					case string:
						f, strconvErr := strconv.ParseFloat(metricValue.(string), 64)
						if strconvErr != nil {
							log.Printf("Unable to convert value for metric: %s\n", strconvErr)
							continue
						}
						metrics[function.ID] = f
						break
					}
				}
			}

		}

	}

	return metrics
}
