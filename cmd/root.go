package cmd

import (
	"fmt"
	"os"

	"github.com/nyambati/aws-service-limits-exporter/internal/metrics"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

const bannerMsg = `
      _    ____   ___  _     _____
     / \  / ___| / _ \| |   | ____|
    / _ \ \___ \| | | | |   |  _|
   / ___ \ ___) | |_| | |___| |___
  /_/   \_\____/ \__\_\_____|_____|
`

func RootCmd() {
	app := cli.NewApp()
	app.Name = "Prometheus service limits exporter"
	app.Usage = "This app exports aws services limits and usage as metric as prometheus metrics"
	app.Description = "This app exports aws services limits and usage as metric as prometheus metrics"
	myFlags := []cli.Flag{
		cli.StringFlag{
			Name:   "listen-address",
			Value:  "0.0.0.0",
			Usage:  "address to listen on for metrics",
			EnvVar: "LISTEN_ADDRESS",
		},
		cli.IntFlag{
			Name:   "port",
			Value:  9090,
			Usage:  "port to listen on for metrics",
			EnvVar: "PORT",
		},
		cli.StringFlag{
			Name:   "config",
			Value:  ".aws-service-limits",
			Usage:  "configuration on how limits should be exported",
			EnvVar: "AWS_SERVICE_LIMITS_CONFIG_FILE",
		},
		cli.StringFlag{
			Name:   "region",
			Value:  "eu-east-1",
			Usage:  "AWS region",
			EnvVar: "AWS_REGION",
		},
	}
	app.Flags = myFlags
	app.Action = run

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func run(c *cli.Context) error {
	fmt.Print(bannerMsg)
	fmt.Println()
	log.Info(fmt.Sprintf("Prometheus exporter starting to listen on %s:%d/metrics", c.String("listen-address"), c.Int("port")))
	metrics.Serve(c)
	return nil
}
