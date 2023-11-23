package metrics

import "github.com/prometheus/client_golang/prometheus"

var registry *prometheus.Registry

type metrics struct {
	acm prometheus.Gauge
}

func Generate(reg prometheus.Registerer) *metrics {
	m := &metrics{
		acm: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "cpu_temperature_celsius",
			Help: "Current temperature of the CPU.",
		}),
	}
	reg.MustRegister(m.acm)
	return m
}

func CreateRegistery() *prometheus.Registry {

	// Create a non-global registry.
	registry = prometheus.NewRegistry()

	// Create new metrics and register them using the custom registry.
	m := Generate(registry)
	// Set values for the new created metrics.
	m.acm.Set(65.3)
	return registry
}
