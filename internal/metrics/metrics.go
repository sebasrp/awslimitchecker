package metrics

import (
	"regexp"
	"strconv"

	"github.com/gobeam/stringy"
	"github.com/nyambati/asqe/internal/services"
	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/urfave/cli"
)

var registry *prometheus.Registry

type ServiceConfig struct {
	Name    string
	Regions []string
}

type Config struct {
	Services             []ServiceConfig
	ServiceQuotaOverride map[string]map[string]float64
}

var config = &Config{}

func CreateMetric(reg prometheus.Registerer, name, help, namespace string) *prometheus.GaugeVec {
	metric := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      name,
			Help:      help,
		},
		[]string{
			"quota_code",
			"resource_id",
			"region",
			"unit",
			"is_global",
		},
	)

	reg.MustRegister(metric)
	return metric
}

func CreateRegistery(context *cli.Context) *prometheus.Registry {
	err := viper.Unmarshal(&config)

	if err != nil {
		log.Fatal(err)
	}
	// Create a non-global registry.
	registry = prometheus.NewRegistry()

	for _, service := range config.Services {
		generateMetrics(service, config.ServiceQuotaOverride)
	}
	return registry
}

func cleanMetricName(quotaName string) string {
	return stringy.New(regexp.MustCompile(`[^a-zA-Z0-9 ]+`).ReplaceAllString(quotaName, " ")).SnakeCase().ToLower()
}

func generateMetrics(service ServiceConfig, overrides map[string]map[string]float64) (metric *prometheus.GaugeVec) {
	quotaOverrides := getQuotaOverrides(overrides)
	for _, region := range service.Regions {
		usage := services.GetUsage(service.Name, region, quotaOverrides)
		for _, quotaInfo := range usage {
			metric = CreateMetric(registry, cleanMetricName(quotaInfo.QuotaName), "", "aws")
			metric.WithLabelValues(quotaInfo.Quotacode, quotaInfo.ResourceId, region, quotaInfo.Unit, strconv.FormatBool(quotaInfo.Global)).Add(quotaInfo.QuotaValue)
		}
	}

	return

}

func getQuotaOverrides(overrides map[string]map[string]float64) []services.AWSQuotaOverride {
	quotaOverrides := []services.AWSQuotaOverride{}
	for svcName, svc := range overrides {
		for quotaName, quota := range svc {
			quotaOverrides = append(quotaOverrides, services.AWSQuotaOverride{Service: svcName, QuotaName: quotaName, QuotaValue: quota})
		}
	}

	return quotaOverrides
}