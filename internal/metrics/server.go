package metrics

import (
	"fmt"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

func Serve(c *cli.Context) {
	http.Handle("/metrics", promhttp.HandlerFor(CreateRegistery(), promhttp.HandlerOpts{}))
	log.Fatal(http.ListenAndServe(fmt.Sprintf("%s:%d", c.String("listen-address"), c.Int("port")), nil))
}
