package metrics

import (
	"fmt"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/urfave/cli"
)

func Serve(context *cli.Context) {
	viper.SetConfigName("config")
	viper.AddConfigPath(context.String("config-path"))
	err := viper.ReadInConfig()

	if err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}

	http.Handle("/metrics", promhttp.HandlerFor(CreateRegistery(context), promhttp.HandlerOpts{}))
	log.Fatal(
		http.ListenAndServe(
			fmt.Sprintf("%s:%d", context.String("listen-address"), context.Int("port")),
			nil,
		),
	)
}
